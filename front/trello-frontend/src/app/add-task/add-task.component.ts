import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Route, Router } from '@angular/router';
import { TaskService } from '../service/task.service';
import { Task } from '../model/task.model';
import { CommonModule } from '@angular/common';


@Component({
  selector: 'app-add-task',
  standalone: true,
  templateUrl: './add-task.component.html',
  styleUrls: ['./add-task.component.css'],
  imports: [ReactiveFormsModule, CommonModule, FormsModule],
  providers: [TaskService]
})
export class AddTaskComponent implements OnInit {
  taskForm: FormGroup;
  projectId: string = '';

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private taskService: TaskService,
    private router: Router,
  ) {
    this.taskForm = this.fb.group({
      title: ['', Validators.required],
      description: ['', Validators.required],
      startDate: ['', Validators.required],
      endDate: ['', Validators.required],
      status: ['Pending', Validators.required],
      projectId: ['', Validators.required]
    });
  }

  ngOnInit(): void {
    // Preuzimanje projectId iz URL-a
    this.route.paramMap.subscribe(params => {
      this.projectId = params.get('projectId') || '';
      // Postavi projectId u formi
      if (this.projectId) {
        this.taskForm.patchValue({ projectId: this.projectId });
      }
    });
  }

  addTask(): void {
    if (this.taskForm.valid) {
      const formValue = this.taskForm.value;

      // Pretvaranje datuma u odgovarajući format (ISO string)
      formValue.startDate = new Date(formValue.startDate).toISOString();
      formValue.endDate = new Date(formValue.endDate).toISOString();

      // Proveri assignedTo
      formValue.assignedTo = formValue.assignedTo ? [formValue.assignedTo] : [];

      const task: Task = formValue;

      this.taskService.createTask(task).subscribe({
        next: (response) => {
          console.log('Task added successfully', response);
          this.taskForm.reset();
          this.router.navigate(['/task-list', this.projectId]);
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
