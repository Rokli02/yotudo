import { Box } from "@mui/material";
import { FC, useMemo, ComponentProps } from "react";
import { SelectedFilename } from "./selected-filename";
import { SelectImage, SelectorType } from "./select-image";

export interface ImageSelectorProps extends Omit<ComponentProps<typeof Box>, 'onChange'> {
    value?: string;
    onChange: (v?: string | null, type?: SelectorType) => void;
    restoreValue: () => void;
}

export const ImageSelector: FC<ImageSelectorProps> = ({ value, onChange, restoreValue, sx, ...props }) => {
    const Content = useMemo(() => {
        switch (value) {
            case null:
            case undefined:
            case "": 
            case "thumbnail":
                return <SelectImage value={value} onChange={onChange} restoreValue={restoreValue} />
            default:
                return <SelectedFilename value={value} onRemove={onChange}/>
        }
    }, [value, onChange])

    return <Box sx={{ position: 'relative', width: '100%', ...sx }} {...props}>
        { Content }
    </Box>;
}
