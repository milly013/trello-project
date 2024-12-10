import { Component, OnInit } from '@angular/core';
import { TaskService } from '../service/task.service';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { threadId } from 'worker_threads';

@Component({
  selector: 'app-task-status',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule,FormsModule],
  templateUrl: './task-status.component.html',
  providers: [TaskService]
})
export class TaskStatusComponent {
  taskId!: string;
  newStatus: string = 'pending';

  constructor(private taskService: TaskService,private route: ActivatedRoute, private router: Router) {}

  updateTaskStatus() {
    this.route.paramMap.subscribe(params=> {
      this.taskId = params.get('id') || '';
      console.log(this.taskId)
    })
    if (this.taskId && this.newStatus) {
      
      this.taskService.updateTaskStatus(this.taskId, this.newStatus).subscribe(
        response => {
          alert('Task status successfully updated');
          // this.router.navigate(['/task-list'])
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