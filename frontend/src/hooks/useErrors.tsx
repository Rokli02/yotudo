import { ReactElement, useMemo, useState } from "react";
import { FormHelperText } from '@mui/material';

export const useError = () => {
    const [errorTexts, setErrorTexts] = useState<FormError>({});
    
    const error = useMemo(() => {
        return Object.entries(errorTexts).reduce((res, [key, value]) => {
            res[key] = (<div>{value.map((v, i) => <FormHelperText key={`${key}_${i}`} error required>{v}</FormHelperText>)}</div>)

            return res;
        }, {} as Record<string, ReactElement>)
    }, [errorTexts]);

    const addError = (key: string, error: string) => {
        setErrorTexts((pre) => {
            if (!pre[key]) pre[key] = [];

            return {
                ...pre,
                [key]: [...pre[key], error]
            }
        })
    }

    const clearErrors = () => {
        setErrorTexts({})
    }

    const clearError = (key: string) => {
        setErrorTexts((pre) => {
            delete pre[key];

            return {
                ...pre
            }
        })
    }

    return {error, addError, clearErrors, clearError}
}

export default useError;

interface FormError extends Record<string, string[]> {}
