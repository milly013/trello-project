import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface User {
  id: string; // ili number, zavisno od vašeg modela
  username: string;
  email: string;
  password: string; // Ako je potrebno
}

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private apiUrl = 'http://localhost:8081'; // URL vašeg API-a

  constructor(private http: HttpClient) {}

  // Metoda za dobijanje liste korisnika
  getUsers(): Observable<User[]> {
    return this.http.get<User[]>(this.apiUrl);
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
}
