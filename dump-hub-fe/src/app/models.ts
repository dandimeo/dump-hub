export interface FileObj {
  filename?: string;
  size?: number;
}

export interface Files {
  files?: FileObj[];
}

export interface SelectedFile {
  file: File;
  uuid: string;
  pending: boolean;
  complete: boolean;
  error: string | null;
  chunkSent: number;
  progress: number;
}
