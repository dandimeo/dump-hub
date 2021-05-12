import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { ApiService } from '../api.service';

interface SearchResponse {
  results?: any[];
  tot?: number;
}

interface PagConfig {
  currentPage: number;
  pageSize: number;
  total: number;
}

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.css'],
})
export class MainComponent implements OnInit {
  searchForm = new FormGroup({
    query: new FormControl(),
  });
  searchError = false;

  results: any[] = [];
  loadingResult = false;
  pagConfig: PagConfig;

  constructor(private apiService: ApiService) {
    this.pagConfig = {
      currentPage: 1,
      pageSize: 20,
      total: 0,
    };
  }

  ngOnInit(): void {
    this.initPaginator();
    this.onQueryChange();
    this.search();
  }

  public search(): void {
    const query = this.searchForm.get('query')?.value;
    this.loadingResult = true;

    this.apiService.search(query, this.pagConfig.currentPage).subscribe(
      (data: SearchResponse) => {
        this.results = [];
        if (data.results && data.tot) {
          this.results = data.results;
          this.pagConfig.total = data.tot;
        }
        this.loadingResult = false;
      },
      (_) => {
        this.results = [];
        this.loadingResult = false;
        this.searchError = true;
      }
    );
  }

  public pageChange(newPage: number): void {
    this.pagConfig.currentPage = newPage;
    this.search();
  }

  private initPaginator(): void {
    this.pagConfig = {
      currentPage: 1,
      pageSize: 20,
      total: 0,
    };
  }

  private onQueryChange(): void {
    this.searchForm.get('query')?.valueChanges.subscribe((_) => {
      this.initPaginator();
      this.search();
    });
  }
}
