import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { UserService } from '../service/user.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-remove-user',
  standalone: true,
  templateUrl: './remove-user.component.html',
  styleUrls: ['./remove-user.component.css'],
  imports: [CommonModule, ReactiveFormsModule],
  providers: [UserService]
})
export class RemoveUserComponent implements OnInit {
  removeUserForm: FormGroup;

  constructor(private fb: FormBuilder, private userService: UserService) {
    this.removeUserForm = this.fb.group({
      projectId: ['', Validators.required],
      userId: ['', Validators.required]
    });
  }

  ngOnInit(): void {}

  removeUser() {
    if (this.removeUserForm.valid) {
      const projectId = this.removeUserForm.get('projectId')?.value;
      const userId = this.removeUserForm.get('userId')?.value;
      this.userService.removeUserFromProject(projectId, userId).subscribe({
        next: (response) => {
          console.log('User removed from project', response);
          this.removeUserForm.reset();
        },
        error: (err) => console.error('Error removing user:', err)
      });
    }
  }
}
