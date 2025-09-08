import { OpenFileDialog as Go_OpenFileDialog, OpenConfirmationDialog as Go_OpenConfirmationDialog } from '@service/DialogService';

export async function OpenFileDialog(): Promise<string> {
    const filename = await Go_OpenFileDialog();
    if (!filename) throw new Error("file not selected");

    return filename;
}

export async function OpenConfirmationDialog(title: string, message: string): Promise<boolean> {
    return await Go_OpenConfirmationDialog(title, message)
}