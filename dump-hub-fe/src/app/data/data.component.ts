import { Component, OnInit } from '@angular/core';
import { ApiService } from '../api.service';
import { PagConfig, Status, StatusData } from '../models';

@Component({
  selector: 'app-data',
  templateUrl: './data.component.html',
  styleUrls: ['./data.component.css'],
})
export class DataComponent implements OnInit {
  loadingStatus = false;
  pagConfig: PagConfig;
  uploads: Status[] = [];
  apiInterval: any;
  deleteModal = false;
  errorMessage = 'Unknown error';
  toDelete: Status | null = null;

  constructor(private apiService: ApiService) {
    this.pagConfig = {
      currentPage: 1,
      pageSize: 20,
      total: 0,
    };
  }

  ngOnInit(): void {
    this.initPaginator();
    this.getStatus();

    this.apiInterval = setInterval(() => {
      this.getStatus();
    }, 5 * 1000);
  }

  ngOnDestroy(): void {
    clearInterval(this.apiInterval);
  }

  public pageChange(newPage: number): void {
    this.pagConfig.currentPage = newPage;
    this.getStatus();
  }

  public onDeleteRequest(item: Status): void {
    this.toDelete = item;
    this.deleteModal = true;
  }

  public onDelete(): void {
    if (this.toDelete) {
      this.deleteHistory(this.toDelete);
    }
    this.deleteModal = false;
  }

  public onDeleteCancel(): void {
    this.toDelete = null;
    this.deleteModal = false;
  }

  public deleteHistory(item: Status): void {
    this.apiService.delete(item.checksum).subscribe(
      (_) => {
        this.toDelete = null;
        this.deleteModal = false;
        this.setDeleteStatus(item);
      },
      (_) => {
        this.errorMessage = 'Unable to delete entries';
      }
    );
  }

  private getStatus(): void {
    this.loadingStatus = true;
    this.apiService.getStatus(this.pagConfig.currentPage).subscribe(
      (data: StatusData) => {
        this.uploads = [];
        this.pagConfig.total = 0;
        if (data.results && data.tot) {
          this.uploads = data.results;
          this.pagConfig.total = data.tot;
        }
        this.loadingStatus = false;
      },
      (_) => {
        this.loadingStatus = true;
      }
    );
  }

  private setDeleteStatus(item: Status): void {
    this.uploads.forEach((upload) => {
      if (upload.checksum === item.checksum) {
        upload.status = 2;
      }
    });
  }

  private initPaginator(): void {
    this.pagConfig = {
      currentPage: 1,
      pageSize: 20,
      total: 0,
    };
  }
}
