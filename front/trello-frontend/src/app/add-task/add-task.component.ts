// add-task.component.ts

import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { TaskService } from '../service/task.service';
import { Task } from '../model/task.model';


@Component({
  selector: 'app-add-task',
  standalone: true,
  templateUrl: './add-task.component.html',
  styleUrls: ['./add-task.component.css'],
  imports: [ReactiveFormsModule],
  providers: [TaskService]
})
export class AddTaskComponent implements OnInit {
  taskForm: FormGroup;

  constructor(private fb: FormBuilder, private taskService: TaskService) {
    this.taskForm = this.fb.group({
      title: ['', Validators.required],
      description: ['', Validators.required],
      startDate: ['', Validators.required],
      endDate: ['', Validators.required],
      assignedTo: ['', Validators.required],
      status: ['Pending', Validators.required],
      projectId: ['', Validators.required]
    });
  }

  ngOnInit(): void {}

  addTask(): void {
    if (this.taskForm.valid) {
      const task: Task = this.taskForm.value;
      
      this.taskService.createTask(task).subscribe({
        next: (response) => {
          console.log('Task added successfully', response);
          this.taskForm.reset();
        },
        error: (error) => {
          console.error('Error adding task', error);
        },
        complete: () => {
          console.log('Add task observable completed');
        }
      });
    }
  }
}
