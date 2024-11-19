import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

interface LoginResponse {
  token: string;
  userId: string;
  userRole: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = 'http://localhost:8000/api/user'; 

  constructor(private http: HttpClient) {}

  login(email: string, password: string): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(`${this.apiUrl}/login`, { email, password });
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
}
