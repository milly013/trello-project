import { Component } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { UserService } from '../service/user.service';


@Component({
  selector: 'app-remove-member-from-task',
  standalone: true,
  templateUrl: './remove-member-from-task.component.html',
  styleUrls: ['./remove-member-from-task.component.css'],
  imports: [ReactiveFormsModule],
  providers: [UserService]
})
export class RemoveMemberFromTaskComponent {
  removeMemberForm: FormGroup; // Reactive form za uklanjanje člana sa taska

  constructor(
    private fb: FormBuilder,
    private userService: UserService // Injektujemo UserService
  ) {
    // Kreiramo formu sa dva inputa: taskId i userId
    this.removeMemberForm = this.fb.group({
      taskId: ['', Validators.required],
      userId: ['', Validators.required]
    });
  }

  // Funkcija koja se poziva pri submitovanju forme
  removeMemberFromTask() {
    if (this.removeMemberForm.valid) {
      const { taskId, userId } = this.removeMemberForm.value;

      // Pozivamo servis koji šalje zahtev za uklanjanje člana sa taska
      this.userService.removeUserFromTask(taskId, userId).subscribe(
        response => {
          console.log('Member removed successfully', response);
          alert('Member removed successfully');
        },
        error => {
          console.error('Error removing member', error);
          alert('Error removing member');
        }
      );
    } else {
      alert('Please fill all the required fields');
    }
  }
}
