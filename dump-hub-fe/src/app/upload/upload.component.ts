import { Component, OnInit } from '@angular/core';
import { ApiService } from '../api.service';
import { FileObj, Files, SelectedFile } from '../models';
import { HttpErrorResponse } from '@angular/common/http';
import * as uuid from 'uuid';
import { concat, Observable } from 'rxjs';

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.css'],
})
export class UploadComponent implements OnInit {
  uploadPath: string = '';
  fileQueue: SelectedFile[] = [];
  files: FileObj[] = [];
  loadingFiles = true;

  constructor(private apiService: ApiService) {}

  ngOnInit(): void {
    this.getFiles();
  }

  private getFiles(): void {
    this.apiService.getFiles().subscribe((data: Files) => {
      this.uploadPath = data.dir;
      this.files = [];
      if (data.files) {
        this.files = data.files;
        this.loadingFiles = false;
      }
    });
  }

  public onSelect(event: any): void {
    if (event.target.files.length > 0) {
      const files = event.target.files;
      for (const inFile of files) {
        const selectedFile: SelectedFile = {
          file: inFile,
          uuid: uuid.v4(),
          pending: true,
          complete: false,
          error: null,
          chunkSent: 0,
          progress: 0,
        };
        this.fileQueue.push(selectedFile);
        this.uploadQueue();
      }
    }
  }

  public uploadQueue(): void {
    this.fileQueue.forEach((selectedFile) => {
      if (selectedFile.pending) {
        selectedFile.pending = false;
        this.uploadFile(selectedFile);
      }
    });
  }

  private uploadFile(selectedFile: SelectedFile): void {
    const chunkSize = 30 * 1000000;
    const fileSize = selectedFile.file.size;
    const chunks = Math.ceil(fileSize / chunkSize);

    let chunk = 0;
    const requests: Observable<any>[] = [];
    while (chunk < chunks) {
      const offset = chunk * chunkSize;
      const slice = selectedFile.file.slice(offset, offset + chunkSize);

      const formData = new FormData();
      formData.append('id', selectedFile.uuid);
      formData.append('filename', selectedFile.file.name);
      formData.append('offset', offset.toString());
      formData.append('file_size', fileSize.toString());
      formData.append('data', slice);

      // const apiResponse = this.apiService.uploadChunk(formData);
      requests.push(this.apiService.uploadChunk(formData));

      /*
      apiResponse.subscribe(
        (_) => {
          selectedFile.chunkSent++;
          selectedFile.progress = Math.ceil(
            (selectedFile.chunkSent * 100) / chunks
          );

          if (selectedFile.chunkSent === chunks) {
            selectedFile.complete = true;
            this.getFiles();
          }
        },
        (err: HttpErrorResponse) => {
          selectedFile.error = 'Unknown error';
          if (typeof err.error === 'string') {
            selectedFile.error = err.error;
          }
          this.getFiles();
          return;
        }
      );
      */

      chunk++;
    }

    concat(...requests).subscribe(
      (_) => {
        selectedFile.chunkSent++;
        selectedFile.progress = Math.ceil(
          (selectedFile.chunkSent * 100) / chunks
        );
        if (selectedFile.chunkSent === chunks) {
          selectedFile.complete = true;
          this.getFiles();
        }
      },
      (err: HttpErrorResponse) => {
        selectedFile.error = 'Unknown error';
        if (typeof err.error === 'string') {
          selectedFile.error = err.error;
        }
        this.getFiles();
        return;
      }
    );
  }
}
