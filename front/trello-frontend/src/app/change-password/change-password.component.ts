import { Component } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { UserService } from '../service/user.service';
import { CommonModule } from '@angular/common';
import { AuthService } from '../service/auth.service';

@Component({
  selector: 'app-change-password',
  standalone: true,
  imports: [CommonModule,ReactiveFormsModule,FormsModule],
  templateUrl: './change-password.component.html',
  styleUrl: './change-password.component.css'
})
export class ChangePasswordComponent {
  changePasswordForm: FormGroup;
  isSubmitting = false;
  errorMessage: string | null = null;

  constructor(
    private formBuilder: FormBuilder,
    private userService: UserService,
    private authService:AuthService,
  ) {
    this.changePasswordForm = this.formBuilder.group({
      currentPassword: ['', [Validators.required]],
      newPassword: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  onSubmit() {
    if (this.changePasswordForm.invalid) {
      return;
    }

    this.isSubmitting = true;
    this.errorMessage = null;

    const { currentPassword, newPassword } = this.changePasswordForm.value;
    const userId = this.authService.getUserId() || '';

    if (!userId) {
      this.errorMessage = 'User not found';
      this.isSubmitting = false;
      return;
    }

    const requestBody = {
      userId,
      currentPassword,
      newPassword
    };

    this.userService.changePassword(requestBody).subscribe(
      response => {
        alert('Password updated successfully');
        this.isSubmitting = false;
      },
      error => {
        this.errorMessage = 'Failed to update password. Please try again.';
        console.error('Error updating password', error);
        this.isSubmitting = false;
      }
    );
  }
}
