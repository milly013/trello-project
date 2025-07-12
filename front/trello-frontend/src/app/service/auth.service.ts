import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import {jwtDecode} from 'jwt-decode';

interface LoginResponse {
  token: string;
  userId: string;
  userRole: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = 'https://localhost/api/user/users'; 
  private tokenKey = 'authToken';

  constructor(private http: HttpClient) {}

  login(email: string, password: string, recaptchaToken: string): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(`${this.apiUrl}/login`, { email, password, recaptchaToken });
  }
  
  // Funkcija za proveru da li je korisnik prijavljen
  isAuthenticated(): boolean {
    return !!localStorage.getItem('authToken');
  }

  // Funkcija za odjavljivanje korisnika
  logout(): void {
    localStorage.removeItem('authToken');
    localStorage.removeItem('managerId')
  }
  // Funkcija za dobijanje headera sa tokenom
  getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('authToken');
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  }
   // Provera da li je korisnik menadžer
   isUserManager(): boolean {
    const role = localStorage.getItem('userRole');
    return role === 'manager';
  }
  
  getUserId(): string | null {
    const token = localStorage.getItem(this.tokenKey);
    if (token) {
      try {
        const decodedToken: any = jwtDecode(token);
        return decodedToken.userId || null;
      } catch (error) {
        console.error('Invalid token', error);
        return null;
      }
    }
    return null;
  }
  getUserRole(): string | null {
    return localStorage.getItem('userRole');
  }

  
  // Provera da li je korisnik član
  isUserMember(): boolean {
    const role = localStorage.getItem('userRole');
    return role === 'member';
  }

  isLoggedIn(): boolean {
    return !!localStorage.getItem('token');
  }
  
}


