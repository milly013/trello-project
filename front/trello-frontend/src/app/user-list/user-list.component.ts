import { Component, OnInit } from '@angular/core';
import { User, UserService } from '../service/user.service';
import { CommonModule } from '@angular/common';
import { AuthService } from '../service/auth.service';

@Component({
  selector: 'app-user-list',
  standalone: true,
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.css'],
  imports:[CommonModule]
})
export class UserListComponent implements OnInit {
  isManager: boolean = false;
  currentUserId: string | null = null;

  users: User[] = [];

  constructor(private userService: UserService, private authService: AuthService) {}

  ngOnInit(): void {
    this.loadUsers();
    this.isManager = this.authService.isUserManager();
    this.currentUserId = this.authService.getUserId();
  }

  loadUsers(): void {
    this.userService.getUsers().subscribe(
      (data) => {
        this.users = data;
      },
      (error) => {
        console.error('Error fetching users:', error);
      }
    );
  }
  deleteUser(userId: string): void {
    this.userService.deleteUser(userId).subscribe(
      () => {
        this.users = this.users.filter(user => user.id !== userId);
      },
      (error) => {
        console.error('Greška pri brisanju korisnika:', error);
      }
    );
  }
}
