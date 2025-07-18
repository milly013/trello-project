import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuthService } from './auth.service';
import { ProjectService } from './project.service';

export interface User {
  id: string;
  username: string;
  email: string;
  password: string; 
  isActive: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private apiUrl = 'https://localhost/api/user/users'; 

  constructor(private http: HttpClient, private authService: AuthService, private projectService: ProjectService) {}

  // Metoda za dobijanje liste korisnika
  getUsers(): Observable<any> {
    const headers = this.authService.getAuthHeaders();
    return this.http.get<any>(`${this.apiUrl}`, { headers });
  }
  getUsersByProjectId(projectId: string): Observable<any>{
    const users = this.projectService.getUsersByProjectId(projectId)
    return users
  }
 
  // Metoda za brisanje korisnika po ID-u
  deleteUser(userId: string): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${userId}`);
  }
  
  // Metoda za registraciju korisnika
  registerUser(userData: { username: string; email: string; password: string }): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });

    return this.http.post<any>(`${this.apiUrl}/register`, userData, { headers });

  }

  // **Nova metoda za slanje verifikacionog koda korisniku** 📧
  sendVerificationCode(email: string): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });

    const body = { email: email };

    return this.http.post<any>(`${this.apiUrl}/send-verification-code`, body, { headers });

  }

  // **Nova metoda za verifikaciju korisnika pomoću koda** ✅
  verifyUser(email: string, code: string): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });

    const body = { email: email, code: code };

    return this.http.post<any>(`${this.apiUrl}/verify`, body, { headers });
  }
  
  forgotPassword(email: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/forgot-password`, { email });
  }

  resetPassword(token: string, newPassword: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/reset-password`, { token, newPassword });
  }

  changePassword(requestBody: { userId: string, currentPassword: string, newPassword: string }): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });
    console.log(requestBody.currentPassword,requestBody.newPassword,requestBody.userId)
    return this.http.post(`${this.apiUrl}/change-password`, requestBody, { headers });

  }
  
  getUserDetails(userId: string): Observable<any> {
    return this.http.get(`${this.apiUrl}/${userId}`);
  }
  
  
}
