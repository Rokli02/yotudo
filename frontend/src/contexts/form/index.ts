import { useContext } from 'react';
import { FormContext } from './FormContext';

export { Form } from './FormContext';
export { FormInput } from './FormInput';
export { FormCheckbox } from './FormCheckbox';
export { FormAutocomplete } from './FormAutocomplete';
export { FormMultiselectAutocomplete } from './FormMultiselectAutocomplete';
export { FormSlider } from './FormSlider';
export const useForm = () => useContext(FormContext)
export type { FormConstraints, AutocompleteOptions } from './interface'