import { CSSObject, styled } from "@mui/material/styles";
import { AuthorService, GenreService, NewMusic } from "@src/api";
import { AutocompleteOptions, FormConstraints } from "@src/contexts/form";
import { AutocompleteProps } from "@src/contexts/form/FormAutocomplete";
import { MultiselectAutocompleteProps } from "@src/contexts/form/FormMultiselectAutocomplete";
import { Dialog, DialogTitle } from "@mui/material";
import { CustomCSS } from "@src/components/common/interface";

export const musicConstraints: FormConstraints = {
    name: (_value, errors) => {
        const value = _value as string;

        if (!value || !value.trim())
            return errors.push('* Kötelező elem')
    },
    published: (_value, errors) => {
        const value = _value as number;

        if (!value)
            return

        if (isNaN(value))
            return errors.push('A mezőnek érvényes számmal írt dátumnak kell lenni')

        const currentYear = new Date().getFullYear();

        if (value < 1900 || value > currentYear)
            return errors.push(`A dátumnak 1900 után, de ${currentYear} előtt kell lenni.`)
    },
    url: (_value, errors) => {
        const value = _value as string;

        if (!value || !value.trim())
            return errors.push('* Kötelező mező')
    },
    author: (_value, errors) => {
        const value = _value as AutocompleteOptions;

        if (!value || typeof value !== 'object' || !value['id'])
            return errors.push('*Kötelezően választani kell egy szerzőt')
    },
    genre: (_value, errors) => {
        const value = _value as AutocompleteOptions;

        if (!value || typeof value !== 'object' || !value['id'])
            return errors.push('*Kötelezően választani kell egy műfajt')
    },
    picName: (_value, errors) => {
        const value = _value as string | undefined | null

        if (value === null)
            return errors.push('Erősítsd meg a választásod, csak utána tudsz menteni!')
    }
}

export const transformFormObjectToNewMusic = (value: Partial<NewMusic>) => {
    value.published = Number(value.published)

    return value;
}

export const getAuthorOptions: AutocompleteProps['getOptions'] = async (search) => {       
    const authors = await AuthorService.GetAuthors({ page: 0, size: 10, filter: search })

    return authors.data.map((author) => {
        const options: AutocompleteOptions = {
            label: author.name,
            ...author,
        }

        return options;
    })
}

export const getGenreOptions: AutocompleteProps['getOptions'] = async (_) => {
    const genres = await GenreService.GetAllGenre();

    return genres.map((genre) => {
        const options: AutocompleteOptions = {
            label: genre.name,
            ...genre,
        }

        return options;
    })
}

//TODO: Átvinni a selectedOptionsId lekezelést backend-be
// export async function getContributorOptions(this: { excludeIds?: number[] }, search: string, selectedOptionsId: number[]): ReturnType<MultiselectAutocompleteProps['getOptions']> {
//     let excludeIds: number[];
//     if (this && 'excludeIds' in this && this.excludeIds?.length) {
//         excludeIds = [...this.excludeIds, ...selectedOptionsId]
//     } else {
//         excludeIds = selectedOptionsId;
//     }
//     const contributors = await AuthorService.GetAuthors({ page: 0, size: 10, filter: search }, excludeIds);
//
//     return contributors.data.filter((c) => !selectedOptionsId.includes(c.id)).map((contributor) => {
//         const options: AutocompleteOptions = {
//             label: contributor.name,
//             ...contributor,
//         }
//
//         return options;
//     })
// }

export const getContributorOptions: MultiselectAutocompleteProps['getOptions'] = async (search: string, selectedOptionsId: number[]) => {
    const contributors = await AuthorService.GetAuthors({ page: 0, size: 10, filter: search });

    return contributors.data.filter((c) => !selectedOptionsId.includes(c.id)).map((contributor) => {
        const options: AutocompleteOptions = {
            label: contributor.name,
            ...contributor,
        }

        return options;
    })
}

export const CustomDialag = styled(Dialog)({
    '& .MuiDialog-container': {
        marginBlock: '4% 6%',
        alignItems: 'flex-start',
        '& .MuiPaper-root': {
            maxWidth: '650px',
            width: '100%',
            backgroundColor: 'var(--background-color)',
            color: 'var(--font-color)',
            '& .MuiDialogContent-root': {
                width: '100%',
                '& > .MuiFormLabel-root': {
                    marginLeft: '6ch',
                },
                '& > .MuiFormControl-root > .MuiSlider-root': {
                    marginLeft: '21px',
                },
            },
        },
    },
    '& .form_items': {
        display: 'flex',
        flexDirection: 'column',
        rowGap: '1rem',
    },
} as CustomCSS)

export const Title = styled(DialogTitle)({
    fontSize: '1.25rem',
    marginInline: 'auto',
    width: 'max-content',
})