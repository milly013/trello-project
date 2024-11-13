<<<<<<< HEAD
import { Component } from '@angular/core';
import { TaskService } from '../service/task.service';
import { FormsModule, } from '@angular/forms';
=======
// add-task.component.ts

import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { TaskService } from '../service/task.service';
import { Task } from '../model/task.model';
import { CommonModule } from '@angular/common';

>>>>>>> develop

@Component({
  selector: 'app-add-task',
  standalone: true,
<<<<<<< HEAD
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
=======
  templateUrl: './add-task.component.html',
  styleUrls: ['./add-task.component.css'],
  imports: [ReactiveFormsModule,CommonModule,FormsModule],
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
>>>>>>> develop
}
