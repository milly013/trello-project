import { Component, OnInit } from '@angular/core';
import { Task } from '../model/task.model';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { TaskService } from '../service/task.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-task-list',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './task-list.component.html',
  styleUrl: './task-list.component.css'
})
export class TaskListComponent implements OnInit {
  projectId!: string;
  tasks: Task[] = [];

  constructor(
    private route: ActivatedRoute,
    private taskService: TaskService
  ) {}

  ngOnInit(): void {
    // Uzimanje projectId iz URL parametara
    this.route.paramMap.subscribe(params => {
      this.projectId = params.get('projectId') || '';
      this.getTasksForProject();
    });
  }

  getTasksForProject(): void {
    if (this.projectId) {
      this.taskService.getTasksByProjectId(this.projectId).subscribe(
        (tasks: Task[]) => {
          this.tasks = tasks;
        },
        error => {
          console.error('Error fetching tasks:', error);
        }
      );
    }
  }
}
