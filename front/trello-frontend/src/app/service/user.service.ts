import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuthService } from './auth.service';

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
  private apiUrl = 'http://localhost:8000/api/user/users'; 

  constructor(private http: HttpClient, private authService: AuthService) {}

  // Metoda za dobijanje liste korisnika
  getUsers(): Observable<any> {
    const headers = this.authService.getAuthHeaders();
    return this.http.get<any>(`${this.apiUrl}`, { headers });
  }
  // Metoda za brisanje korisnika po ID-u
  deleteUser(userId: string): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${userId}`);
  }


  // Metoda za dodavanje novog korisnika u projekat
  addUserToProject(projectId: string, userId: string): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });

    const body = {
      userId: userId
    };

    return this.http.post<any>(`${this.apiUrl}/projects/${projectId}/members`, body, { headers });
  }

  // Metoda za uklanjanje korisnika iz projekta
  removeUserFromProject(projectId: string, userId: string): Observable<any> {
    const url = `${this.apiUrl}/projects/${projectId}/members`;
    const options = {
      headers: {
        'Content-Type': 'application/json'
      },
      body: { userId: userId }
    };
    return this.http.delete(url, options);
  }

  // Metoda za registraciju korisnika
  registerUser(userData: { username: string; email: string; password: string }): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });

    return this.http.post<any>(`${this.apiUrl}/users/register`, userData, { headers });
  }

  // **Nova metoda za slanje verifikacionog koda korisniku** 📧
  sendVerificationCode(email: string): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });

    const body = { email: email };

    return this.http.post<any>(`${this.apiUrl}/users/send-verification-code`, body, { headers });
  }

  // **Nova metoda za verifikaciju korisnika pomoću koda** ✅
  verifyUser(email: string, code: string): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });

    const body = { email: email, code: code };

    return this.http.post<any>(`${this.apiUrl}/users/verify`, body, { headers });
  }

  addMemberToTask(taskId: string, userId: string): Observable<any> {
    const url = `${this.apiUrl}/tasks/${taskId}/members`;
    const body = { userId: userId };
    return this.http.post<any>(url, body, {
      headers: new HttpHeaders({ 'Content-Type': 'application/json' })
    });
  }


  removeUserFromTask(taskId: string, userId: string): Observable<any> {
    const url = `${this.apiUrl}/tasks/${taskId}/members/${userId}`;
    return this.http.delete(url);
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
    return this.http.post(`${this.apiUrl}/change-password`, requestBody, { headers });
  }
  
}
