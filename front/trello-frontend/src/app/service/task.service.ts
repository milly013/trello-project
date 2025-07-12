
// service/task.service.ts

import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Task } from '../model/task.model';

@Injectable({
  providedIn: 'root'
})
export class TaskService {
  private apiUrl = 'https://localhost/api/task'; 

  constructor(private http: HttpClient) {}

  createTask(task: Task): Observable<Task> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });
    return this.http.post<Task>(`${this.apiUrl}/tasks`, task, { headers });

  }
  getTaskById(taskId: string): Observable<Task>{
    return this.http.get<Task>(`${this.apiUrl}/task/${taskId}`);
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
  addMemberToTask(taskId: string, userId: string): Observable<any> {
    const url = `${this.apiUrl}/tasks/add-member`;
    const body = { 
      userId: userId,
      taskId: taskId
     };
    return this.http.post<any>(url, body, {
      headers: new HttpHeaders({ 'Content-Type': 'application/json' })
    });
  }
  removeUserFromTask(taskId: string, userId: string): Observable<any> {
    const url = `${this.apiUrl}/tasks/remove-member`;
    const body = {
      userId: userId,
      taskId: taskId
    };
  
    // Koristimo `request` metod da šaljemo DELETE zahtev sa telom
    return this.http.request<any>('delete', url, {
      headers: new HttpHeaders({ 'Content-Type': 'application/json' }),
      body: body
    });
  }
  
}
