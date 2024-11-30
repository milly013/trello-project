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
  getProjectsByManager(managerId: string): Observable<Project[]> {
    return this.http.get<Project[]>(`${this.apiUrl}/projects/manager/${managerId}`);
  }
  getProjectsByMember(memberId: string): Observable<Project[]>{
    return this.http.get<Project[]>(`${this.apiUrl}/projects/member/${memberId}`);
  }

  getProjectById(id: string): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    })
    return this.http.get<any>(`${this.apiUrl}/projects/${id}`, { headers });
  }
  getUsersByProjectId(projectId: string): Observable<any> {
    const url = `${this.apiUrl}/users/${projectId}`;
    return this.http.get<any>(url, {
      headers: new HttpHeaders({ 'Content-Type': 'application/json' })
    });
  }
  addUserToProject(projectId: string, memberId: string): Observable<any>{
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });
    const body = {
      memberId: memberId
    };
    return this.http.post<any>(`${this.apiUrl}/projects/${projectId}/members`, body, { headers });
  }
  
}
