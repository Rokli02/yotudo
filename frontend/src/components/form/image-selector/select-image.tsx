import { Box, IconButton, SxProps, Theme } from "@mui/material";
import { FC, ReactNode, useMemo, useState } from "react";
import { Option, Select, SelectProps } from "../select";
import { Button, FormControl, Input, InputLabel } from "..";
import { Done, Replay } from "@mui/icons-material";
import { DialogService } from '@src/api'

export type SelectorType = 'none' | 'thumbnail' | 'web' | 'local'

export interface SelectImageProps {
    value?: "" | "thumbnail";
    onChange: (type: SelectorType, v?: string | null) => void;
    onError?: (e: string) => void;
    restoreValue: () => void;
}

export const SelectImage: FC<SelectImageProps> = ({ value, onChange, restoreValue }) => {
    const [option, setOption] = useState<Option<SelectorType>>(SelectOptions['none']);

    const Content: ReactNode = useMemo(() => {
        switch (option.value) {
            case 'none':
            case 'thumbnail': {
                return <>
                    <div />
                    <IconButton sx={ContentStyle.ActionButton} onClick={restoreValue} title="Visszaállítás">
                        <Replay />
                    </IconButton>
                </>
            }
            case 'web': {
                return <SelectImageWeb onChange={onChange}/>
            }
            case 'local': {
                return <Button
                    variant="outlined"
                    color="secondary"
                    onClick={() => {
                        DialogService.OpenFileDialog()
                            .then((selectedFilename) => onChange('local', selectedFilename))
                            .catch(() => {})
                    }}
                >Kiválasztás</Button>
            }
        }
    }, [option, value])

    const onSelectType: SelectProps['onChange'] = (chosenValue: SelectorType) => {
        switch (chosenValue) {
            case 'none':
                onChange('none', undefined);

                break;
            case 'thumbnail':
                onChange('thumbnail', chosenValue);

                break;
            default:
                onChange('none', null);
        }

        setOption(SelectOptions[chosenValue])
    }

    return <Box sx={ContainerStyle}>
        <Select
            sx={ContentStyle.Select}
            options={SelectOptionsArray}
            value={option.value}
            onChange={onSelectType}
        />
        { Content }
    </Box>
}

const SelectOptions: Record<SelectorType, Option<SelectorType>> = {
    'none': {
        label: 'Nincs',
        value: 'none',
    },
    'thumbnail': {
        label: 'Borító',
        value: 'thumbnail',
    },
    'web': {
        label: 'Web',
        value: 'web',
    },
    'local': {
        label: 'Helyi',
        value: 'local',
    },
}
const SelectOptionsArray = Object.values(SelectOptions)


const ContainerStyle: SxProps<Theme> = {
    display: 'inline-grid',
    alignItems: 'center',
    width: '100%',
    gridTemplateColumns: '100px auto min-content',
    columnGap: '10px'
}

const ContentStyle = {
    Select: {
        height: '46px',
    },
    ActionButton: {
        color: 'var(--font-color)',
        padding: '4px',
        width: 30,
        height: 30,
        '>svg': {
            width: '100%',
            height: '100%',
        }
    },
    WebFormControl: {
        alignSelf: 'start',
        display: 'flex',
        alignItems: 'center',
        height: 'min-content'
    },
} as const satisfies Record<string, SxProps<Theme>>

const SelectImageWeb: FC<{ onChange: (type: SelectorType, v?: string | undefined) => void }> = ({ onChange }) => {
    const [state, setState] = useState<string>('');

    const confirm = () => onChange('web', state)

    return <>
        <FormControl sx={ContentStyle.WebFormControl} variant="standard">
            <InputLabel>URL</InputLabel>
            <Input
                onChange={(e) => {
                    e.preventDefault();
                    
                    setState(e.target.value);
                }}
                
                />
        </FormControl>
        <IconButton sx={ContentStyle.ActionButton} onClick={confirm} title="Megerősít">
            <Done />
        </IconButton>
    </>
}