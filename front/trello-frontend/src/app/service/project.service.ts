import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Project {
  name: string;
  expected_end_date: string;
  min_members: number;
  max_members: number;
  manager_id: string;
}

@Injectable({
  providedIn: 'root'
})
export class ProjectService {

  private apiUrl = 'http://localhost:8081'; 
  
  constructor(private http: HttpClient) { }

  createProject(project: Project): Observable<Project> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });
    return this.http.post<Project>(`${this.apiUrl}/projects`, project, { headers });
  }
}
