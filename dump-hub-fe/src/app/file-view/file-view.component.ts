import { Component, Input, OnInit } from '@angular/core';
import { ApiService } from '../api.service';
import { FileObj, Files } from '../models';

@Component({
  selector: 'app-file-view',
  templateUrl: './file-view.component.html',
  styleUrls: ['./file-view.component.css'],
})
export class FileViewComponent implements OnInit {
  @Input() loadingFiles: boolean;
  @Input() uploadPath: string;
  @Input() fileError: string | null;
  @Input() files: FileObj[];

  toDelete: FileObj | undefined = undefined;
  loadingResult = true;
  deleteModal = false;

  constructor(private apiService: ApiService) {
    this.loadingFiles = true;
    this.uploadPath = '';
    this.fileError = null;
    this.files = [];
  }

  ngOnInit(): void {}

  private getFiles(): void {
    this.apiService.getFiles().subscribe(
      (data: Files) => {
        this.files = [];
        this.uploadPath = data.dir;
        if (data.files) {
          this.files = data.files;
          this.fileError = null;
          this.loadingFiles = false;
        }
      },
      (_) => {
        this.fileError = 'Unable to retrieve folder content';
      }
    );
  }

  private deleteFile(toDelete: FileObj): void {
    if (toDelete.filename) {
      const id = btoa(toDelete.filename);
      this.apiService.deleteFile(id).subscribe(
        (_) => {
          this.loadingFiles = true;
          this.getFiles();

          this.fileError = null;
          this.toDelete = undefined;
          this.deleteModal = false;
        },
        (_) => {
          this.fileError = 'Unable to delete selected file';
          this.toDelete = undefined;
          this.deleteModal = false;
        }
      );
    }
  }

  public openDelete(file: FileObj | undefined): void {
    if (file) {
      this.toDelete = file;
      this.deleteModal = true;
    }
  }

  public onDelete(): void {
    if (this.toDelete) {
      this.deleteFile(this.toDelete);
    }
  }

  public onDeleteCancel(): void {
    this.toDelete = undefined;
    this.deleteModal = false;
  }

  public errorModal(): boolean {
    if (this.fileError != null) {
      return true;
    }

    return false;
  }
}
