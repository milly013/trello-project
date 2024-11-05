import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class TaskService {
  private apiUrl = 'http://localhost:8082/tasks';

  constructor(private http: HttpClient) {}

  addTask(task: any): Observable<any> {
    return this.http.post<any>(this.apiUrl, task);
  }
}
