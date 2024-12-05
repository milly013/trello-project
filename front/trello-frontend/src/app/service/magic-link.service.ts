import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class MagicLinkService {
  private apiUrl = 'https://localhost/api/user/users'; // URL backend API-ja

  constructor(private http: HttpClient) {}

  // Method to request magic link
  requestMagicLink(email: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/request-magic-link`, { email });
  }

  // Method to verify magic link token
  verifyMagicLinkToken(token: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/magic-login`, { token });
  }
}