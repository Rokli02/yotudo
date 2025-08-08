import { CSSObject } from "@mui/material/styles";

export interface CustomCSS extends CSSObject, Record<`&${string}` | `@${string}`, CustomCSS> {};