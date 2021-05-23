import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ApiService } from '../api.service';
import { Alert, FileObj, Files, Preview } from '../models';

@Component({
  selector: 'app-analyze',
  templateUrl: './analyze.component.html',
  styleUrls: ['./analyze.component.css'],
})
export class AnalyzeComponent implements OnInit {
  files: FileObj[] = [];
  preview: string[] = [];
  previewTable: string[][] = [];
  previewTableMaxCols = 0;
  selectedFile: FileObj | null;
  loadingFiles = true;
  history: Alert[] = [];

  patternForm = new FormGroup({
    startLine: new FormControl('', Validators.required),
    separator: new FormControl('', Validators.required),
  });

  analyzeForm = new FormGroup({
    pattern: new FormControl('', Validators.required),
    columns: new FormControl('', Validators.required),
  });

  constructor(private apiService: ApiService) {
    this.selectedFile = null;
    this.patternForm.setValue({
      startLine: '0',
      separator: ':',
    });
    this.preview = ['Select a file to enable preview...'];
  }

  ngOnInit(): void {
    this.getFiles();
    this.onStartChange();
    this.onSepChange();
  }

  public analyzeFile(): void {
    if (this.selectedFile) {
      if (this.selectedFile.filename) {
        let fn = this.selectedFile.filename;
        let pattern = this.analyzeForm.get('pattern')?.value;
        let columns = this.analyzeForm.get('columns')?.value;
        this.apiService.analyze(fn, pattern, columns).subscribe(
          (_) => {
            this.history.push({
              message: fn + ' will be analyzed in background',
              type: 0,
            });

            this.files.forEach((f, index) => {
              if (f == this.selectFile) delete this.files[index];
            });

            this.previewTableMaxCols = 0;
            this.previewTable = [];
            this.selectedFile = null;
            this.preview = [];
            this.analyzeForm.get('pattern')?.setValue(null);
            this.analyzeForm.get('columns')?.setValue(null);

            document.body.scrollTop;
          },
          (_) => {
            this.history.push({
              message: 'unable to analyze file: ' + fn,
              type: -1,
            });
            this.loadingFiles = true;
            this.getFiles();

            this.previewTableMaxCols = 0;
            this.previewTable = [];
            this.selectedFile = null;
            this.preview = [];
            this.analyzeForm.get('pattern')?.setValue(null);
            this.analyzeForm.get('columns')?.setValue(null);

            document.body.scrollTop;
          }
        );
      }
    }
  }

  public selectFile(file: FileObj): void {
    if (file != this.selectedFile) {
      this.selectedFile = file;
      this.getPreview(file);
    }
  }

  public isSelected(file: FileObj): boolean {
    if (file == this.selectedFile) {
      return true;
    }
    return false;
  }

  public isColumnSelected(colNumber: number): boolean {
    const selected: number[] = this.analyzeForm.get('columns')?.value;
    if (selected.indexOf(colNumber) === -1) {
      return false;
    }
    return true;
  }

  public toggleColumn(colNumber: number): void {
    const selected = this.analyzeForm.get('columns')?.value;
    if (this.isColumnSelected(colNumber)) {
      const index: number = selected.indexOf(colNumber);
      if (index !== -1) {
        selected.splice(index, 1);
      }

      this.analyzeForm.get('columns')?.setValue(selected);
      return;
    }

    selected.push(colNumber);
    this.analyzeForm.get('columns')?.setValue(selected);
  }

  public counter(i: number): Array<number> {
    return new Array(i);
  }

  public removeHistory(alert: Alert) {
    this.history.forEach((element, index) => {
      if (element == alert) {
        delete this.history[index];
      }
    });
  }

  private onStartChange(): void {
    this.patternForm.get('startLine')?.valueChanges.subscribe((_) => {
      if (this.selectedFile) {
        this.getPreview(this.selectedFile);
      }
    });
  }

  private onSepChange(): void {
    this.patternForm.get('separator')?.valueChanges.subscribe((_) => {
      if (this.selectedFile && this.preview.length) {
        this.parsePreviewTable();
      }
    });
  }

  private getPreview(file: FileObj): void {
    this.preview = ['Loading preview data...'];
    const start = parseInt(this.patternForm.get('startLine')?.value);
    if (file.filename) {
      this.apiService.preview(file.filename, start).subscribe(
        (data: Preview) => {
          this.preview = data.preview;
          this.parsePreviewTable();
        },
        () => {
          this.preview = ['Unable to get file preview...'];
        }
      );
    }
  }

  private parsePreviewTable(): void {
    this.previewTable = [];
    this.previewTableMaxCols = 0;

    const separator = this.patternForm.get('separator')?.value;
    this.preview.forEach((content) => {
      const values = content.replace(' ', '').split(separator);
      if (values.length > this.previewTableMaxCols) {
        this.previewTableMaxCols = values.length;
      }
    });

    this.preview.forEach((content) => {
      const tableRow: string[] = [];
      for (let i = 0; i < this.previewTableMaxCols; i++) {
        tableRow[i] = 'N/A';
      }
      const values = content.replace(' ', '').split(separator);
      for (let j = 0; j < values.length; j++) {
        if (values[j].length > 1) {
          tableRow[j] = values[j];
        }
      }
      this.previewTable.push(tableRow);
    });

    this.analyzeForm.get('columns')?.setValue([]);
    this.computePattern();
  }

  private getFiles(): void {
    this.apiService.getFiles().subscribe((data: Files) => {
      this.files = [];
      if (data.files) {
        this.files = data.files;
        this.loadingFiles = false;
        if (this.files.length) {
          this.selectFile(this.files[0]);
        }
      }
    });
  }

  private computePattern(): void {
    const start = this.patternForm.get('startLine')?.value;
    const sep = this.patternForm.get('separator')?.value;
    let pattern = '{' + start + '}';
    pattern = pattern + '{' + sep + '}';
    this.analyzeForm.get('pattern')?.setValue(pattern);
  }
}
