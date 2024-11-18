import { Component } from '@angular/core';
import { TaskService } from '../service/task.service';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-task-status',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule,FormsModule],
  templateUrl: './task-status.component.html',
  providers: [TaskService]
})
export class TaskStatusComponent {
  taskId: string = '';
  newStatus: string = 'pending';

  constructor(private taskService: TaskService) {}

  updateTaskStatus() {
    if (this.taskId && this.newStatus) {
      this.taskService.updateTaskStatus(this.taskId, this.newStatus).subscribe(
        response => {
          alert('Task status successfully updated');
        },
        error => {
          console.error('Error updating task status', error);
          alert('An error occurred while updating the task status');
        }
      );
    } else {
      alert('Please enter all required fields.');
    }
  }
}