import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private UPLOAD = environment.baseAPI + 'upload';
  private SEARCH = environment.baseAPI + 'search';
  private HISTORY = environment.baseAPI + 'history';
  private DELETE = environment.baseAPI + 'delete';
  private FILES = environment.baseAPI + 'files';

  constructor(private httpClient: HttpClient) {}

  public uploadChunk(data: FormData): Observable<any> {
    return this.httpClient.post(this.UPLOAD, data);
  }

  public search(query: string, page: number): Observable<any> {
    const data = {
      query,
      page,
    };
    return this.httpClient.post(this.SEARCH, data);
  }

  public getHistory(page: number): Observable<any> {
    const data = {
      page,
    };
    return this.httpClient.post(this.HISTORY, data);
  }

  public delete(checkSum: string): Observable<any> {
    const data = {
      checkSum,
    };
    return this.httpClient.post(this.DELETE, data);
  }

  public getFiles(): Observable<any> {
    return this.httpClient.get(this.FILES);
  }

  public deleteFile(id: string): Observable<any> {
    return this.httpClient.delete(this.FILES + '/' + id);
  }
}
