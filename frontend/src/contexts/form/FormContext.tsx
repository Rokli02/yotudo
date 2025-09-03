/* eslint-disable @typescript-eslint/no-explicit-any */
import { createContext, Dispatch, FC, FormEvent, SetStateAction, useRef } from "react";
import { Form as DefaultForm } from "../../components/form";
import useError from "../../hooks/useErrors";
import { IForm, IFormContext } from "./interface";

export const FormContext = createContext<IFormContext | null>(null);

export const Form: FC<IForm> = ({
    children,
    onSubmit,
    transformFlatObjectTo = (v: object) => v,
    FormComponent = DefaultForm,
    clearOnSubmit = true,
    constraints,
}) => {
    const waitingForResponse = useRef<boolean>(false)
    const inputs = useRef<Inputs>({});
    const { error, addError, clearErrors, clearError } = useError();

    const registerInput: IFormContext['registerInput'] = (name, setValue, onClear, initValue) => {
        inputs.current[name] = { value: initValue, setValue, onClear, initValue };
    }

    const unregisterInput: IFormContext['unregisterInput'] = (name) => {
        delete inputs.current[name];
        clearError(name);
    }

    const onValueChange: IFormContext['onValueChange'] = (name: string, value: unknown) => {
        inputs.current[name]['value'] = value;
    }

    const _onSubmit = async (event: FormEvent) => {
        event.preventDefault();

        if (waitingForResponse.current) return false

        waitingForResponse.current = true
        let hasError = false

        try {          
            const flatObject = Object.entries(inputs.current).reduce((res, [key, { value }]) => {
                res[key] = value;
                
                return res;
            }, {} as Record<string, Inputs['']['value']>);
            
            if (constraints) {
                clearErrors();
                
                for (const [name, constraint] of Object.entries(constraints)) {
                    const errors: string[] = [];
                    constraint(flatObject[name], errors, flatObject);
                    
                    if (!hasError && errors.length > 0) hasError = true
                    
                    errors.forEach((error) => addError(name, error));
                }
            }
            
            if (hasError) return;

            const resOfSubmit = await onSubmit(transformFlatObjectTo(flatObject));
                        
            if (resOfSubmit && clearOnSubmit) {
                clearInputs(inputs.current);
            }
        } catch (err) {
            console.error('Some error occured in "FormContext":', err)
        } finally {
            waitingForResponse.current = false
        }
    }

    const getErrors: IFormContext['getErrors'] = (name) => error[name];

    return (
        <FormContext.Provider value={{
            onValueChange,
            registerInput,
            unregisterInput,
            getErrors,
        }}>
            <FormComponent onSubmit={_onSubmit}>
                {children}
            </FormComponent>
        </FormContext.Provider>
    )
}

interface Inputs extends Record<string, {
    value: any;
    setValue: Dispatch<SetStateAction<any>>;
    onClear: () => void;
    initValue?: unknown;
}> {}

const clearInputs = (inputs: Inputs) => {
    Object.values(inputs).forEach((input) => {
        input.value = input.initValue;
        input.onClear()
    })
}
