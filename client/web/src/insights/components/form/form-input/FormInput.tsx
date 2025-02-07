import classnames from 'classnames'
import React, { useRef, forwardRef, InputHTMLAttributes, ReactNode, useEffect } from 'react'
import { useMergeRefs } from 'use-callback-ref'

import { LoaderInput } from '@sourcegraph/branded/src/components/LoaderInput'

import styles from './FormInput.module.scss'
import { ForwardReferenceComponent } from './types'

interface FormInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'title'> {
    /** Title of input. */
    title?: ReactNode
    /** Description block for field. */
    description?: ReactNode
    /** Custom class name for root label element. */
    className?: string
    /** Error massage for input. */
    error?: string
    /** Prop to control error input element state. */
    errorInputState?: boolean
    /** Valid sign to show valid state on input. */
    valid?: boolean
    /** Turn on loading state (visually this is an input with loader) */
    loading?: boolean
    /** Turn on or turn off autofocus for input. */
    autofocus?: boolean
    /** Custom class name for input element. */
    inputClassName?: string
    /** Input icon (symbol) which render right after the input element. */
    inputSymbol?: ReactNode
}

/**
 * Displays the input with description, error message, visual invalid and valid states.
 */
const FormInput = forwardRef((props, reference) => {
    const {
        as: Component = 'input',
        type = 'text',
        title,
        description,
        className,
        inputClassName,
        inputSymbol,
        valid,
        error,
        loading = false,
        errorInputState,
        autoFocus,
        ...otherProps
    } = props

    const localReference = useRef<HTMLInputElement>(null)
    const mergedReference = useMergeRefs([localReference, reference])

    useEffect(() => {
        if (autoFocus) {
            // In some cases if form input has been rendered within reach/portal element
            // react autoFocus set focus too early and in this case we have to
            // call focus explicitly in the next tick to be sure that focus will be
            // on input element. See reach/portal implementation and notice async way to
            // render children in react portal component.
            // https://github.com/reach/reach-ui/blob/0ae833201cf842fc00859612cfc6c30a593d593d/packages/portal/src/index.tsx#L45
            requestAnimationFrame(() => {
                localReference.current?.focus()
            })
        }
    }, [autoFocus])

    return (
        <label className={classnames('w-100', className)}>
            {title && <div className="mb-2">{title}</div>}

            <LoaderInput className="d-flex" loading={loading}>
                <Component
                    type={type}
                    className={classnames(styles.input, inputClassName, 'form-control', 'with-invalid-icon', {
                        'is-valid': valid,
                        'is-invalid': !!error || errorInputState,
                    })}
                    {...otherProps}
                    autoFocus={autoFocus}
                    ref={mergedReference}
                />

                {inputSymbol}
            </LoaderInput>

            {error && (
                <small role="alert" className="text-danger form-text">
                    {error}
                </small>
            )}
            {!error && description && (
                <small className={classnames('text-muted', 'form-text', styles.description)}>{description}</small>
            )}
        </label>
    )
}) as ForwardReferenceComponent<'input', FormInputProps>

FormInput.displayName = 'FormInput'

export { FormInput }
