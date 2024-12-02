import { Component, OnInit } from '@angular/core';
import { NotificationService } from '../service/notification.service';
import { AuthService } from '../service/auth.service';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Notification } from '../model/notification.model'; // Koristi model iz jednog fajla

@Component({
  selector: 'app-user-notifications',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, FormsModule],
  templateUrl: './user-notifications.component.html',
  styleUrls: ['./user-notifications.component.css']
})
export class UserNotificationsComponent implements OnInit{
  notifications: Notification[] = [];
  userId: string = '';

  constructor(private notificationService: NotificationService, private authService: AuthService) { }

  ngOnInit(): void {
    this.userId = this.authService.getUserId() || '';
    if (this.userId) {
      this.getNotifications();
    }
  }

  getNotifications(): void {
    this.notificationService.getNotificationsByUserId(this.userId).subscribe(
      (notifications: Notification[]) => {
        this.notifications = notifications;
        console.log(this.notifications)
      },
      (error) => {
        console.error('Error fetching notifications:', error);
      }
    );
  }
}
