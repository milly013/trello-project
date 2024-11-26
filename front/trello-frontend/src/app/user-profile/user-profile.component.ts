import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { UserService } from '../service/user.service';
import { AuthService } from '../service/auth.service';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-user-profile',
  standalone: true,
  imports: [CommonModule,ReactiveFormsModule,FormsModule,RouterLink ],
  templateUrl: './user-profile.component.html',
  styleUrl: './user-profile.component.css'
})
export class UserProfileComponent implements OnInit {
  user: any;
  errorMessage: string | null = null;

  constructor(
    private userService: UserService,
    private authService: AuthService
  ) {}

  ngOnInit() {
    const userId = this.authService.getUserId();
    if (userId) {
      this.userService.getUserDetails(userId).subscribe(
        response => {
          this.user = response;
        },
        error => {
          console.error('Error fetching user details', error);
          this.errorMessage = 'Failed to load user details. Please try again.';
        }
      );
    } else {
      this.errorMessage = 'User not found';
    }
  }
}