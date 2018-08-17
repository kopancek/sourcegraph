import { DiffPart, JumpURLFetcher } from '@sourcegraph/codeintellify'
import { Controller } from 'cxp/module/environment/controller'
import { Extension } from 'cxp/module/environment/extension'
import { ConfigurationCascade, TextDocumentPositionParams } from 'cxp/module/protocol'
import { HoverMerged } from 'cxp/module/types/hover'
import { Observable, of, OperatorFunction, throwError as error } from 'rxjs'
import { ajax, AjaxResponse } from 'rxjs/ajax'
import { catchError, map, tap } from 'rxjs/operators'
import { Definition, TextDocumentIdentifier } from 'vscode-languageserver-types'
import { DidOpenTextDocumentParams, InitializeResult, ServerCapabilities } from 'vscode-languageserver/lib/main'
import {
    AbsoluteRepo,
    AbsoluteRepoFile,
    AbsoluteRepoFilePosition,
    AbsoluteRepoLanguageFile,
    FileSpec,
    makeRepoURI,
    parseRepoURI,
    PositionSpec,
    RepoSpec,
    ResolvedRevSpec,
    RevSpec,
} from '../repo'
import { getModeFromPath, isPrivateRepository, repoUrlCache, sourcegraphUrl, supportedModes } from '../util/context'
import { memoizeObservable } from '../util/memoize'
import { toAbsoluteBlobURL } from '../util/url'
import { normalizeAjaxError } from './errors'
import { getHeaders } from './headers'

export interface LSPRequest {
    method: string
    params: any
}

/** LSP proxy error code for unsupported modes */
export const EMODENOTFOUND = -32000

export function isEmptyHover(hover: HoverMerged | null): boolean {
    return !hover || !hover.contents || (Array.isArray(hover.contents) && hover.contents.length === 0)
}

function wrapLSP(req: LSPRequest, ctx: AbsoluteRepo, path: string): any[] {
    return [
        {
            id: 0,
            method: 'initialize',
            params: {
                rootUri: `git://${ctx.repoPath}?${ctx.commitID}`,
                initializationOptions: { mode: `${getModeFromPath(path)}` },
            },
        },
        {
            id: 1,
            ...req,
        },
        {
            id: 2,
            method: 'shutdown',
        },
        {
            // id not included on 'exit' requests
            method: 'exit',
        },
    ]
}

/**
 * Inspects a response from LSP Proxy and throws an exception if the response
 * has an error. This is intended to be used in rxjs: pipe(...throwIfError)
 */
const extractLSPResponse: OperatorFunction<AjaxResponse, any> = source =>
    source.pipe(
        tap(ajaxResponse => {
            // Workaround for https://github.com/ReactiveX/rxjs/issues/3606
            if (ajaxResponse.status === 0) {
                throw Object.assign(new Error('Ajax status 0'), ajaxResponse)
            }
        }),
        catchError(err => {
            normalizeAjaxError(err)
            throw err
        }),
        map<AjaxResponse, any[]>(({ response }) => response),
        tap(lspResponses => {
            for (const lspResponse of lspResponses) {
                if (lspResponse && lspResponse.error) {
                    throw Object.assign(new Error(lspResponse.error.message), lspResponse.error, {
                        responses: lspResponses,
                    })
                }
            }
        }),
        map(lspResponses => lspResponses[1] && lspResponses[1].result)
    )

const fetchHover = memoizeObservable((pos: AbsoluteRepoFilePosition): Observable<HoverMerged | null> => {
    const mode = getModeFromPath(pos.filePath)
    if (!mode || !supportedModes.has(mode)) {
        return of({ contents: [] })
    }

    const body = wrapLSP(
        {
            method: 'textDocument/hover',
            params: {
                textDocument: {
                    uri: `git://${pos.repoPath}?${pos.commitID}#${pos.filePath}`,
                },
                position: {
                    character: pos.position.character! - 1,
                    line: pos.position.line - 1,
                },
            },
        },
        pos,
        pos.filePath
    )

    const url = repoUrlCache[pos.repoPath] || sourcegraphUrl
    if (!url) {
        throw new Error('Error fetching hover: No URL found.')
    }
    if (!canFetchForURL(url)) {
        return of(null)
    }

    return ajax({
        method: 'POST',
        url: `${url}/.api/xlang/textDocument/hover`,
        headers: getHeaders(),
        crossDomain: true,
        withCredentials: true,
        body: JSON.stringify(body),
        async: true,
    }).pipe(extractLSPResponse)
}, makeRepoURI)

const fetchDefinition = memoizeObservable((pos: AbsoluteRepoFilePosition): Observable<Definition> => {
    const mode = getModeFromPath(pos.filePath)
    if (!mode || !supportedModes.has(mode)) {
        return of([])
    }

    const body = wrapLSP(
        {
            method: 'textDocument/definition',
            params: {
                textDocument: {
                    uri: `git://${pos.repoPath}?${pos.commitID}#${pos.filePath}`,
                },
                position: {
                    character: pos.position.character! - 1,
                    line: pos.position.line - 1,
                },
            },
        },
        pos,
        pos.filePath
    )

    const url = repoUrlCache[pos.repoPath] || sourcegraphUrl
    if (!url) {
        throw new Error('Error fetching definition: No URL found.')
    }
    if (!canFetchForURL(url)) {
        return of([])
    }

    return ajax({
        method: 'POST',
        url: `${url}/.api/xlang/textDocument/definition`,
        headers: getHeaders(),
        crossDomain: true,
        withCredentials: true,
        body: JSON.stringify(body),
        async: true,
    }).pipe(extractLSPResponse)
}, makeRepoURI)

export function fetchJumpURL(
    fetchDefinition: SimpleCXPFns['fetchDefinition'],
    pos: AbsoluteRepoFilePosition
): Observable<string | null> {
    return fetchDefinition(pos).pipe(
        map(def => {
            const defArray = Array.isArray(def) ? def : [def]
            def = defArray[0]
            if (!def) {
                return null
            }

            const uri = parseRepoURI(def.uri) as AbsoluteRepoFilePosition
            uri.position = { line: def.range.start.line + 1, character: def.range.start.character + 1 }
            return toAbsoluteBlobURL(uri)
        })
    )
}

export type JumpURLLocation = RepoSpec & RevSpec & ResolvedRevSpec & FileSpec & PositionSpec & { part?: DiffPart }
export function createJumpURLFetcher(
    fetchDefinition: SimpleCXPFns['fetchDefinition'],
    buildURL: (pos: JumpURLLocation) => string
): JumpURLFetcher {
    return ({ line, character, part, commitID, repoPath, ...rest }) =>
        fetchDefinition({ ...rest, commitID, repoPath, position: { line, character } }).pipe(
            map(def => {
                const defArray = Array.isArray(def) ? def : [def]
                def = defArray[0]
                if (!def) {
                    return null
                }

                const uri = parseRepoURI(def.uri)
                return buildURL({
                    repoPath: uri.repoPath,
                    commitID: uri.commitID!, // LSP proxy always includes a commitID in the URI.
                    rev: uri.repoPath === repoPath && uri.commitID === commitID ? rest.rev : uri.rev!, // If the commitID is the same, keep the rev.
                    filePath: uri.filePath!, // There's never going to be a definition without a file.
                    position: {
                        line: def.range.start.line + 1,
                        character: def.range.start.character + 1,
                    },
                    part,
                })
            })
        )
}

/**
 * Modes that are known to not be supported because the server replied with a mode not found error
 */
const unsupportedModes = new Set<string>()

const fetchServerCapabilities = (pos: AbsoluteRepoLanguageFile): Observable<ServerCapabilities | undefined> => {
    // Check if mode is known to not be supported
    const mode = getModeFromPath(pos.filePath)
    if (!mode || unsupportedModes.has(mode)) {
        return error(Object.assign(new Error('Language not supported'), { code: EMODENOTFOUND }))
    }

    const body = wrapLSP(
        {
            method: 'textDocument/didOpen',
            params: {
                textDocument: {
                    uri: `git://${pos.repoPath}?${pos.commitID}#${pos.filePath}`,
                },
            } as DidOpenTextDocumentParams,
        },
        pos,
        pos.filePath
    )
    const url = repoUrlCache[pos.repoPath] || sourcegraphUrl
    if (!url) {
        throw new Error('Error fetching server capabilities. No URL found.')
    }
    if (!canFetchForURL(url)) {
        return of(undefined)
    }
    return ajax({
        method: 'POST',
        url: `${url}/.api/xlang/textDocument/didOpen`,
        headers: getHeaders(),
        crossDomain: true,
        withCredentials: true,
        body: JSON.stringify(body),
        async: true,
    }).pipe(
        tap(response => {
            if (response.status === 0) {
                throw Object.assign(new Error('Ajax status 0'), response)
            }
        }),
        catchError<AjaxResponse, never>(err => {
            normalizeAjaxError(err)
            throw err
        }),
        map(({ response }) => response),
        map(results => {
            for (const result of results) {
                if (result && result.error) {
                    if (result.error.code === EMODENOTFOUND) {
                        unsupportedModes.add(mode)
                    }
                    throw Object.assign(new Error(result.error.message), result.error)
                }
            }

            return results.map((result: any) => result && result.result)
        }),
        map(results => (results[0] as InitializeResult).capabilities)
    )
}

function canFetchForURL(url: string): boolean {
    if (url === 'https://sourcegraph.com' && isPrivateRepository()) {
        return false
    }
    return true
}

export interface SimpleCXPFns {
    fetchHover: (pos: AbsoluteRepoFilePosition) => Observable<HoverMerged | null>
    fetchDefinition: (pos: AbsoluteRepoFilePosition) => Observable<Definition>
    fetchServerCapabilities: (pos: AbsoluteRepoLanguageFile) => Observable<ServerCapabilities | undefined>
}

export const lspViaAPIXlang: SimpleCXPFns = {
    fetchHover,
    fetchDefinition,
    fetchServerCapabilities,
}

export const toTextDocumentIdentifier = (pos: AbsoluteRepoFile): TextDocumentIdentifier => ({
    uri: `git://${pos.repoPath}?${pos.commitID}#${pos.filePath}`,
})

const toTextDocumentPositionParams = (pos: AbsoluteRepoFilePosition): TextDocumentPositionParams => ({
    textDocument: toTextDocumentIdentifier(pos),
    position: {
        character: pos.position.character! - 1,
        line: pos.position.line - 1,
    },
})

export const createLSPViaCXP = <X extends Extension, C extends ConfigurationCascade>(
    cxpController: Controller<X, C>
) => ({
    fetchHover: pos =>
        cxpController.registries.textDocumentHover
            .getHover(toTextDocumentPositionParams(pos))
            .pipe(map(hover => (hover === null ? HoverMerged.from([]) : hover))),
    fetchDefinition: pos =>
        cxpController.registries.textDocumentDefinition.getLocation(toTextDocumentPositionParams(pos)),
    fetchServerCapabilities,
})
