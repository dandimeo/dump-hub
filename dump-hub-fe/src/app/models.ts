export interface FileObj {
  filename?: string;
  size?: number;
}

export interface Files {
  dir: string;
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

export interface Preview {
  preview: string[];
}

export interface Alert {
  type: number;
  message: string;
}

export interface PagConfig {
  currentPage: number;
  pageSize: number;
  total: number;
}

export interface Status {
  date: string;
  filename: string;
  checksum: string;
  status: number;
}

export interface StatusData {
  results?: Status[];
  tot?: number;
}
