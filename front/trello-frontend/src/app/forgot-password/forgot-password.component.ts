import { Component } from '@angular/core';
import { UserService } from '../service/user.service';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-forgot-password',
  standalone: true,
  imports: [CommonModule,ReactiveFormsModule,FormsModule],
  templateUrl: './forgot-password.component.html',
  styleUrl: './forgot-password.component.css'
})
export class ForgotPasswordComponent {
  email: string = '';

  constructor(private userService: UserService) {}

  forgotPassword() {
    if (this.email) {
      this.userService.forgotPassword(this.email).subscribe(
        response => {
          alert('Password reset email sent successfully');
        },
        error => {
          console.error('Error sending reset email', error);
          alert('An error occurred while sending the reset email');
        }
      );
    } else {
      alert('Please enter your email address.');
    }
  }
}