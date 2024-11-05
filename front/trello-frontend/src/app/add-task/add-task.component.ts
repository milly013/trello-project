import { Component } from '@angular/core';
import { TaskService } from '../service/task.service';
import { FormsModule, } from '@angular/forms';

@Component({
  selector: 'app-add-task',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './add-task.component.html',
})
export class AddTaskComponent {
  task = {
    title: '',
    description: '',
    projectId: '',
    startDate: '',
    endDate: '',
    assignedTo: '',
    status: 'Pending'
  };

  constructor(private taskService: TaskService) {}

  onSubmit() {
    this.taskService.addTask(this.task).subscribe({
      next: (response) => {
        console.log('Task successfully added:', response);
        alert('Task successfully added!');
      },
      error: (error) => {
        console.error('Error adding task:', error);
        alert('Error adding task');
      },
      complete: () => {
        console.log('Task addition completed.');
      }
    });
  }
}
