/*
The MIT License (MIT)
Copyright (c) 2021 Davide Pataracchia
Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:
The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.
*/

import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.css']
})
export class UploadComponent implements OnInit {
  uploadForm = new FormGroup({
    pattern: new FormControl('', Validators.required),
    file: new FormControl('', Validators.required),
    columns: new FormControl('', Validators.required),
  });

  patternForm = new FormGroup({
    separator: new FormControl('', Validators.required),
    commentChar: new FormControl('', Validators.required)
  });

  uploadStatus = 0;
  editPatternModal = false;

  fileContent: string[] = [];
  fileContentRaw: any;
  selectedFile: any;

  previewContent: string[] = [];
  previewTable: string[][] = [];
  previewTableMaxCols = 0;

  constructor(
    private apiService: ApiService
  ) { }

  ngOnInit(): void {
    this.patternForm.setValue({
      separator: ':',
      commentChar: '#'
    });
    this.patternString();

    this.uploadForm.controls.pattern.disable();
    this.uploadForm.get('file')?.setValue(null);
    this.uploadForm.get('columns')?.setValue([]);
    this.previewContent = ['Select a text file to enable preview...']

    this.onPatternChange();
    this.onCommentChange();
  }

  public onSubmit(): void {
    const formData = new FormData();
    formData.append('file', this.uploadForm.get('file')?.value);
    formData.append('pattern', this.uploadForm.get('pattern')?.value);
    formData.append('columns', this.uploadForm.get('columns')?.value);
    this.uploadStatus = 1;

    this.apiService.upload(formData)
      .subscribe(
        (_) => {
          this.uploadStatus = 2;
        },
        _ => this.uploadStatus = -1
      );
  }

  public onFileSelect(event: any): void {
    this.uploadForm.controls.file.setValue(null);
    this.uploadForm.controls.columns.setValue([]);
    this.previewContent = ['Loading file...'];

    if (event.target.files.length > 0) {
      const file = event.target.files[0];
      if (!file.type.startsWith('text/')) {
        this.previewContent = ['Invalid file type.'];
        this.uploadForm.controls.file.setValue(null);
        return;
      }

      this.uploadForm.controls.file.setValue(file);
      this.readFile();
    }
  }

  public patternString(): void {
    const separator = this.patternForm.get('separator')?.value;
    const commentChar = this.patternForm.get('commentChar')?.value;

    const value = `{${separator}}{${commentChar}}`;
    this.uploadForm.controls.pattern.setValue(value);
  }

  public loadingModal(): boolean {
    return this.uploadStatus !== 0;
  }

  public counter(i: number): Array<number> {
    return new Array(i);
  }

  public isSelected(colNumber: number): boolean {
    const selected: number[] = this.uploadForm.get('columns')?.value;
    if (selected.indexOf(colNumber) === -1) {
      return false;
    }
    return true;
  }

  public toggleColumn(colNumber: number): void {
    const selected = this.uploadForm.get('columns')?.value;
    if (this.isSelected(colNumber)) {
      const index: number = selected.indexOf(colNumber);
      if (index !== -1) {
        selected.splice(index, 1);
      }

      this.uploadForm.get('columns')?.setValue(selected);
      return;
    }

    selected.push(colNumber)
    this.uploadForm.get('columns')?.setValue(selected);
  }

  private readFile(): void {
    this.fileContentRaw = null;
    const file = this.uploadForm.get('file')?.value;
    if (file == null) {
      return;
    }

    const reader: FileReader = new FileReader();
    reader.readAsText(file.slice(0, 8192));
    reader.onloadend = () => {
      this.fileContentRaw = reader.result;
      this.processPreview();
    };

    reader.onerror = () => {
      this.uploadForm.controls.file.setValue(null);
      this.previewContent = ['Unable to read the input file.'];
      return;
    };
  }

  private parsePreview(): void {
    if (this.uploadForm.get('file')?.value == null) {
      return;
    }
    const separator = this.patternForm.get('separator')?.value;

    this.previewTable = [];
    this.previewTableMaxCols = 0;
    this.previewContent.forEach(content => {
      const values = content.replace(' ', '').split(separator);
      if (values.length > this.previewTableMaxCols) {
        this.previewTableMaxCols = values.length;
      }
    });

    this.previewContent.forEach(content => {
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

    this.uploadForm.get('columns')?.setValue([]);
  }

  private processPreview(): void {
    this.fileContent = this.fileContentRaw.split(/[\r\n]+/g);

    this.previewContent = [];

    for (const content of this.fileContent) {
      if (this.previewContent.length >= 20) {
        break;
      }

      const commentChar = this.patternForm.get('commentChar')?.value;
      if (content.replace(' ', '').charAt(0) === commentChar) {
        continue;
      }
      this.previewContent.push(content);
    }

    this.parsePreview();
  }

  private onPatternChange(): void {
    this.uploadForm.get('pattern')?.valueChanges
      .subscribe(_ => {
        this.parsePreview();
      });
  }

  private onCommentChange(): void {
    this.patternForm.get('commentChar')?.valueChanges
      .subscribe(_ => {
        this.processPreview();
      });
  }
}
