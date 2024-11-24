
// service/task.service.ts

import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Task } from '../model/task.model';

@Injectable({
  providedIn: 'root'
})
export class TaskService {
  private apiUrl = 'http://localhost:8000/api/task'; 

  constructor(private http: HttpClient) {}

  createTask(task: Task): Observable<Task> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });
    return this.http.post<Task>(`${this.apiUrl}/tasks`, task, { headers });

  }

  updateTaskStatus(taskId: string, newStatus: string): Observable<any> {
    return this.http.put(`${this.apiUrl}/tasks/${taskId}/status`, { new_status: newStatus });
  }
  getTasksByProjectId(projectId: string): Observable<Task[]> {
    return this.http.get<Task[]>(`${this.apiUrl}/tasks/${projectId}/tasks`);
  }
  getUsersByTaskId(taskId: string): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/tasks/members/${taskId}/users`);
  }
  
}
