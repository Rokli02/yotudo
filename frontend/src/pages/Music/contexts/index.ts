import { useContext } from "react";
import { MusicContext } from "./MusicContext";

export const useMusicContext = () => useContext(MusicContext);
export { MusicProvider } from './MusicContext';