import { OpenFileDialog as Go_OpenFileDialog } from '@service/DialogService';

export async function OpenFileDialog(): Promise<string> {
    const filename = await Go_OpenFileDialog();
    if (!filename) throw new Error("file not selected");

    return filename;
}