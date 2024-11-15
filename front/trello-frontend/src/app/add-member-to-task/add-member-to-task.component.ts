import { Component } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { UserService } from '../service/user.service';


@Component({
  selector: 'app-add-member-to-task',
  standalone: true,
  templateUrl: './add-member-to-task.component.html',
  styleUrls: ['./add-member-to-task.component.css'],
  imports: [ReactiveFormsModule],
  providers: [UserService]
})
export class AddMemberToTaskComponent {
  memberForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private userService: UserService
  ) {
    // Kreiramo Reactive Form
    this.memberForm = this.fb.group({
      taskId: ['', Validators.required],
      userId: ['', Validators.required]
    });
  }

  // Funkcija koja poziva UserService za dodavanje člana na task
  addMemberToTask() {
    if (this.memberForm.valid) {
      const taskId = this.memberForm.get('taskId')?.value;
      const userId = this.memberForm.get('userId')?.value;

      this.userService.addMemberToTask(taskId, userId).subscribe(
        response => {
          console.log('Member added to task:', response);
          alert('Member added to task successfully!');
        },
        error => {
          console.error('Error adding member to task:', error);
          alert('Failed to add member to task.');
        }
      );
    }
  }
}
