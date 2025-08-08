import { useContext } from "react";
import { AuthorContext } from "./AuthorContext";

export const useAuthorContext = () => useContext(AuthorContext);
export { AuthorProvider } from './AuthorContext';
