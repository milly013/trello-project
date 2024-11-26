import { Component } from '@angular/core';
import { UserService } from '../service/user.service';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-reset-password',
  standalone: true,
  imports: [CommonModule,ReactiveFormsModule,FormsModule],
  templateUrl: './reset-password.component.html',
  styleUrl: './reset-password.component.css'
})
export class ResetPasswordComponent {
  token: string = '';
  newPassword: string = '';

  constructor(private userService: UserService) {}

  resetPassword() {
    if (this.token && this.newPassword) {
      this.userService.resetPassword(this.token, this.newPassword).subscribe(
        response => {
          alert('Password updated successfully');
        },
        error => {
          console.error('Error resetting password', error);
          alert('An error occurred while resetting the password');
        }
      );
    } else {
      alert('Please enter all required fields.');
    }
  }
}