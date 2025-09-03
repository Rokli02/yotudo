/* eslint-disable @typescript-eslint/no-explicit-any */
import { DetailedHTMLProps, Dispatch, FC, ReactElement, SetStateAction } from "react";

export interface IForm {
    children: ReactElement | ReactElement[],
    onSubmit: (value: any) => Promise<boolean>,
    transformFlatObjectTo?: (value: any) => any,
    FormComponent?: FC<DetailedHTMLProps<React.FormHTMLAttributes<HTMLFormElement>, HTMLFormElement>>,
    clearOnSubmit?: boolean,
    constraints?: FormConstraints,
}

export interface IFormContext {
    onValueChange: (name: string, value: unknown) => void;
    registerInput: (name: string, setValue: Dispatch<SetStateAction<any>>, onClear: () => void, initValue?: unknown) => void;
    unregisterInput: (name: string) => void;
    getErrors: (name: string) => ReactElement;
}

export interface FormConstraints extends Record<string, (value: unknown, errors: string[], values: Record<string, any>) => void> {}

export interface AutocompleteOptions {
    id?: number;
    label: string;
    [k: string]: unknown;
}