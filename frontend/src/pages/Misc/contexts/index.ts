import { useContext } from "react";
import { GenreContext } from "./GenreContext";

export const useGenreContext = () => useContext(GenreContext);
export { GenreProvider } from './GenreContext';
