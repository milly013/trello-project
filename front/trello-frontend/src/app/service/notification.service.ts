import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Notification } from '../model/notification.model'; // Koristi isti model

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  private baseUrl = 'https://localhost/api/notification/notifications';

  constructor(private http: HttpClient) { }

  getNotificationsByUserId(userId: string): Observable<Notification[]> {
    return this.http.get<Notification[]>(`${this.baseUrl}/${userId}`);
  }
}
