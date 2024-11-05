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

  // Metoda za dodavanje novog korisnika
  addUserToProject(projectId: string, userId: string): Observable<any> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });
  
    // Definišemo telo zahteva
    const body = {
      userId: userId // ID korisnika se šalje u telu zahteva
    };
  
    // Priprema API poziva sa ID projekta u putanji
    return this.http.post<any>(`${this.apiUrl}/projects/${projectId}/members`, body, { headers });
  }
  

  // Metoda za uklanjanje korisnika
  removeUser(userId: string): Observable<void> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });
    return this.http.delete<void>(`${this.apiUrl}/${userId}`, { headers });
  }
}
