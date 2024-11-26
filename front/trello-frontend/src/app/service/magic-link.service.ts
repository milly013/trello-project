import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class MagicLinkService {
  private apiUrl = 'http://localhost:8000/api/user'; // URL backend API-ja

  constructor(private http: HttpClient) {}

  // Metoda za slanje magic link zahteva
  requestMagicLink(email: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/request-magic-link`, { email });
  }

  // Metoda za verifikaciju magic link tokena
  verifyMagicLinkToken(token: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/magic-login`, { token });
  }
}