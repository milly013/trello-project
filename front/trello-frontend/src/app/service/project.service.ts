import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Project {
  id: string;
  name: string;
  endDate: Date;
  minMembers: number;
  maxMembers: number;
  managerId: string;
  isActive: boolean;
  createdAt: Date;
  memberIds: string[];
  taskIds: string[];  
}

@Injectable({
  providedIn: 'root'
})
export class ProjectService {
  private apiUrl = 'http://localhost:8000/api/project'; 
  
  constructor(private http: HttpClient) { }

  createProject(project: Project): Observable<Project> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });
    return this.http.post<Project>(`${this.apiUrl}/projects`, project, { headers });
  }

  getProjects(): Observable<Project[]> {
    return this.http.get<Project[]>(`${this.apiUrl}/projects`);
  }

  getProjectById(id: string): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    })
    return this.http.get<any>(`${this.apiUrl}/projects/${id}`, { headers });
  }
}
