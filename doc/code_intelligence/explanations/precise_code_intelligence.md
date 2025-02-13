# Precise code intelligence

Precise code intelligence relies on [LSIF](https://github.com/Microsoft/language-server-protocol/blob/master/indexFormat/specification.md) 
(Language Server Index Format) data to deliver precomputed code intelligence. It provides fast and highly accurate code intelligence but needs to be periodically generated and uploaded to your Sourcegraph instance. Precise code intelligence is an opt-in feature: repositories for which you have not uploaded LSIF data will continue to use the search-based code intelligence.

> NOTE: Precise code intelligence using LSIF is supported in Sourcegraph 3.8 and up.

## Getting started

See the [how-to guides](../how-to/index.md) to get started with precise code intelligence.

## Cross-repository code intelligence

Cross-repository code intelligence works out-of-the-box when both the dependent repository and the dependency repository has LSIF data _at the correct commits or versions_. We are working on relaxing this constraint so that nearest-commit functionality works on a cross-repository basis as well.

When the current repository has LSIF data and a dependent doesn't, the missing precise results will be supplemented with imprecise search-based code intelligence. This also applies when both repositories have LSIF data, but for a different set of versions. For example, if repository A@v1 depends on B@v2, then we will get precise cross-repository intelligence when we have LSIF data for both A@v1 and B@v2, but would not get a precise result we instead have LISF data for A@v1 and B@v1.

## Why are my results sometimes incorrect?

If LSIF data is not found for a particular file in a repository, Sourcegraph will fall back to search-based code intelligence. You may occasionally see results from [search-based code intelligence](search_based_code_intelligence.md) even when you have uploaded LSIF data. Such results are indicated with a ![tooltip](../img/basic-code-intel-tooltip.svg) tooltip. This can happen in the following scenarios:

- The symbol has LSIF data, but it is defined in a repository which does not have LSIF data.
- The nearest commit that has LSIF data is too far away from your browsing commit. [The limit is 100 commits](https://github.com/sourcegraph/sourcegraph/blob/e7803474dbac8021e93ae2af930269045aece079/lsif/src/shared/constants.ts#L25) ahead/behind.
- The line containing the symbol was created or edited between the nearest indexed commit and the commit being browsed.
- The _Find references_ panel will always include search-based results, but only after all of the precise results have been displayed. This ensures every symbol has code intelligence.

## Size of upload data

The following table gives a rough estimate for the space and time requirements for indexing and conversion. These repositories are a representative sample of public Go repositories available on GitHub. The working tree size is the size of the clone at the given commit (without git history), the number of files indexed, and the number of lines of Go code in the repository. The index size gives the size of the uncompressed LSIF output of the indexer. The conversion size gives the total amount of disk space occupied after uploading the dump to a Sourcegraph instance.

| Repository | Working tree size | Index time | Index size | Processing time | Post-processing size |
| ------------------------------------------------------------------- | ------------------------------- | ------ | ----- | ------- | ----- |
| [bigcache](https://github.com/allegro/bigcache/tree/b7689f7)        | 216KB,   32 files,   2.585k loc |  1.18s | 3.5MB |   0.45s | 0.6MB |
| [sqlc](https://github.com/kyleconroy/sqlc/tree/16cc4e9)             | 396KB,   24 files,   7.041k loc |  1.53s | 7.2MB |   1.62s | 1.6MB |
| [nebula](https://github.com/slackhq/nebula/tree/a680ac2)            | 700KB,   71 files,  10.704k loc |  2.48s |  16MB |   1.63s | 2.9MB |
| [cayley](https://github.com/cayleygraph/cayley/tree/4d89b8a)        | 5.6MB,  226 files,  36.346k loc |  5.58s |  51MB |   4.68s |  11MB |
| [go-ethereum](https://github.com/ethereum/go-ethereum/tree/275cd49) |  27MB,  945 files, 317.664k loc | 20.53s | 255MB |  77.40s |  50MB |
| [kubernetes](https://github.com/kubernetes/kubernetes/tree/e680ad7) | 301MB, 4577 files,   1.550m loc |  1.21m | 910MB |  80.06s | 162MB |
| [aws-sdk-go](https://github.com/aws/aws-sdk-go/tree/18a2d30)        | 119MB, 1759 files,   1.067m loc |  8.20m | 1.3GB | 155.82s | 358MB |

## Data retention policy

The bulk of LSIF data is stored on-disk, and as code intelligence data for a commit ages it becomes less useful. Sourcegraph will automatically remove the least recently uploaded data if the amount of used disk space exceeds a configurable threshold. This value defaults to 10 GiB (10⨉2^30 = 10737418240  bytes), and can be changed via the `DBS_DIR_MAXIMUM_SIZE_BYTES` environment variable.

## More about LSIF

- [Writing an LSIF indexer](writing_an_indexer.md)
- [Adding LSIF to many repositories](../how-to/adding_lsif_to_many_repos.md)

To learn more, check out our lightning talk about LSIF from GopherCon 2019 or the [introductory blog post](https://about.sourcegraph.com/go/code-intelligence-with-lsif):

<iframe width="560" height="315" src="https://www.youtube.com/embed/fMIRKRj_A88" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
