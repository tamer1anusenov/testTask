import React, { useState, useEffect } from 'react';
import { Plus, Search, Filter, SortAsc, List, Calendar as CalendarIcon } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Calendar } from '@/components/ui/calendar';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';
import { format } from 'date-fns';
import { TaskItem } from './TaskItem';
import { ThemeToggle } from './ThemeToggle';
import { WeekView } from './WeekView';
import { cn } from '@/lib/utils';
import { CreateTask, GetAllTasks, DeleteTask, ToggleTaskStatus, UpdateTask } from '../../wailsjs/go/main/App';
import { models } from '../../wailsjs/go/models';

export interface Task {
  id: string;
  title: string;
  completed: boolean;
  createdAt: Date;
  deadline: Date;
  priority: 'low' | 'medium' | 'high';
  archived: boolean;
}

export type FilterType = 'all' | 'active' | 'completed' | 'archived';
export type PriorityType = 'all' | 'low' | 'medium' | 'high';
export type SortType = 'newest' | 'oldest' | 'alphabetical';
export type ViewType = 'list' | 'calendar';

export const TodoApp: React.FC = () => {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [newTaskTitle, setNewTaskTitle] = useState('');
  const [newTaskDeadline, setNewTaskDeadline] = useState<Date>(new Date());
  const [newTaskPriority, setNewTaskPriority] = useState<'low' | 'medium' | 'high'>('medium');
  const [filter, setFilter] = useState<FilterType>('all');
  const [priorityFilter, setPriorityFilter] = useState<PriorityType>('all');
  const [sort, setSort] = useState<SortType>('newest');
  const [view, setView] = useState<ViewType>('list');
  const [currentWeek, setCurrentWeek] = useState(new Date());
  
  // State for delete confirmation modal
  const [taskToDelete, setTaskToDelete] = useState<Task | null>(null);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

  // Helper function to map backend models.Task to frontend Task
  const mapBackendTask = (backendTask: models.Task): Task => ({
    id: backendTask.id.toString(),
    title: backendTask.title,
    completed: backendTask.status === 'completed',
    createdAt: new Date(backendTask.created_at),
    deadline: backendTask.due_date ? new Date(backendTask.due_date) : new Date(),
    priority: backendTask.priority as 'low' | 'medium' | 'high',
    archived: backendTask.archived || false,
  });

  // Load tasks from backend on component mount
  useEffect(() => {
    const loadTasks = async () => {
      try {
        const backendTasks = await GetAllTasks();
        const mappedTasks = backendTasks.map(mapBackendTask);
        setTasks(mappedTasks);
      } catch (error) {
        console.error('Failed to load tasks from backend:', error);
        // Set empty array if backend fails
        setTasks([]);
      }
    };
    
    loadTasks();
  }, []);

    const addTask = async () => {
    if (newTaskTitle.trim()) {
      try {
        const deadlineStr = newTaskDeadline.toISOString().split('T')[0]; // Format as YYYY-MM-DD
        const createdTask = await CreateTask(newTaskTitle.trim(), '', newTaskPriority, deadlineStr);
        const newTask = mapBackendTask(createdTask);
        setTasks(prev => [newTask, ...prev]);
        setNewTaskTitle('');
        setNewTaskDeadline(new Date()); // Reset to today
        setNewTaskPriority('medium'); // Reset to default
      } catch (error) {
        console.error('Failed to create task:', error);
        // Fallback to local creation if backend fails
        const newTask: Task = {
          id: Date.now().toString(),
          title: newTaskTitle.trim(),
          completed: false,
          createdAt: new Date(),
          deadline: newTaskDeadline,
          priority: newTaskPriority,
          archived: false,
        };
        setTasks(prev => [newTask, ...prev]);
        setNewTaskTitle('');
        setNewTaskDeadline(new Date());
        setNewTaskPriority('medium');
      }
    }
  };

  const toggleTask = async (id: string) => {
    try {
      const updatedTask = await ToggleTaskStatus(parseInt(id));
      const mappedTask = mapBackendTask(updatedTask);
      setTasks(prev =>
        prev.map(task =>
          task.id === id ? mappedTask : task
        )
      );
    } catch (error) {
      console.error('Failed to toggle task status:', error);
      // Fallback to local update if backend fails
      setTasks(prev =>
        prev.map(task =>
          task.id === id ? { ...task, completed: !task.completed } : task
        )
      );
    }
  };

  // Function to initiate delete confirmation
  const initiateDeleteTask = (taskId: string) => {
    const task = tasks.find(t => t.id === taskId);
    if (task) {
      setTaskToDelete(task);
      setIsDeleteDialogOpen(true);
    }
  };

  // Function to actually delete the task after confirmation
  const confirmDeleteTask = async () => {
    if (!taskToDelete) return;
    
    try {
      await DeleteTask(parseInt(taskToDelete.id));
      setTasks(prev => prev.filter(task => task.id !== taskToDelete.id));
    } catch (error) {
      console.error('Failed to delete task:', error);
      // Fallback to local deletion if backend fails
      setTasks(prev => prev.filter(task => task.id !== taskToDelete.id));
    } finally {
      setIsDeleteDialogOpen(false);
      setTaskToDelete(null);
    }
  };

  // Function to cancel delete confirmation
  const cancelDeleteTask = () => {
    setIsDeleteDialogOpen(false);
    setTaskToDelete(null);
  };

  // Function to archive a completed task
  const archiveTask = async (taskId: string) => {
    try {
      // For now, simulate the archive API call
      // await ArchiveTask(parseInt(taskId));
      
      // Update task locally to set archived = true
      setTasks(prev =>
        prev.map(task =>
          task.id === taskId ? { ...task, archived: true } : task
        )
      );
    } catch (error) {
      console.error('Failed to archive task:', error);
      // Fallback to local update if backend fails
      setTasks(prev =>
        prev.map(task =>
          task.id === taskId ? { ...task, archived: true } : task
        )
      );
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      addTask();
    }
  };

  // Filter tasks
  const filteredTasks = tasks.filter(task => {
    // Filter by status/archive
    let statusMatch = true;
    switch (filter) {
      case 'active':
        statusMatch = !task.completed && !task.archived;
        break;
      case 'completed':
        statusMatch = task.completed && !task.archived;
        break;
      case 'archived':
        statusMatch = task.archived;
        break;
      default: // 'all'
        statusMatch = true; // Show all tasks including archived
    }

    // Filter by priority
    const priorityMatch = priorityFilter === 'all' || task.priority === priorityFilter;

    return statusMatch && priorityMatch;
  });

  // Filter tasks for week view (exclude archived tasks)
  const weekViewTasks = tasks.filter(task => {
    // Exclude archived tasks from week view
    if (task.archived) return false;
    
    // Apply same filtering as list view but exclude archived
    let statusMatch = true;
    switch (filter) {
      case 'active':
        statusMatch = !task.completed;
        break;
      case 'completed':
        statusMatch = task.completed;
        break;
      case 'archived':
        statusMatch = false; // Never show archived in week view
        break;
      default: // 'all'
        statusMatch = true; // Show all non-archived tasks
    }

    // Filter by priority
    const priorityMatch = priorityFilter === 'all' || task.priority === priorityFilter;

    return statusMatch && priorityMatch;
  });

  // Sort tasks for list view
  const sortedTasks = [...filteredTasks].sort((a, b) => {
    switch (sort) {
      case 'oldest':
        return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
      case 'alphabetical':
        return a.title.localeCompare(b.title);
      case 'newest':
      default:
        return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
    }
  });

  // Sort tasks for week view
  const sortedWeekTasks = [...weekViewTasks].sort((a, b) => {
    switch (sort) {
      case 'oldest':
        return new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
      case 'alphabetical':
        return a.title.localeCompare(b.title);
      case 'newest':
      default:
        return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
    }
  });

  const activeCount = tasks.filter(task => !task.completed && !task.archived).length;
  const completedCount = tasks.filter(task => task.completed && !task.archived).length;

  return (
    <div className="min-h-screen bg-app-bg">
      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Header */}
        <header className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-4xl font-bold bg-gradient-primary bg-clip-text text-transparent">
              My To-Do List
            </h1>
            <p className="text-muted-foreground mt-2">
              {activeCount} active, {completedCount} completed
            </p>
          </div>
          <ThemeToggle />
        </header>

        {/* Add Task Section */}
        <div className="bg-card rounded-lg p-6 shadow-md mb-8 border border-app-border">
          <div className="space-y-4">
            <div className="flex gap-3">
              <Input
                placeholder="Add a new task..."
                value={newTaskTitle}
                onChange={(e) => setNewTaskTitle(e.target.value)}
                onKeyPress={handleKeyPress}
                className="flex-1 bg-input border-app-border focus:border-primary"
              />
              <Select value={newTaskPriority} onValueChange={(value: 'low' | 'medium' | 'high') => setNewTaskPriority(value)}>
                <SelectTrigger className="w-32 bg-input border-app-border">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="low">🟢 Low</SelectItem>
                  <SelectItem value="medium">🟡 Medium</SelectItem>
                  <SelectItem value="high">🔴 High</SelectItem>
                </SelectContent>
              </Select>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    className={cn(
                      "w-48 justify-start text-left font-normal border-app-border hover:bg-app-surface-hover",
                      "bg-input"
                    )}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4" />
                    {format(newTaskDeadline, "MMM d, yyyy")}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={newTaskDeadline}
                    onSelect={(date) => date && setNewTaskDeadline(date)}
                    initialFocus
                    className="p-3 pointer-events-auto"
                  />
                </PopoverContent>
              </Popover>
              <Button 
                onClick={addTask} 
                className="bg-gradient-primary hover:opacity-90 transition-opacity"
                disabled={!newTaskTitle.trim()}
              >
                <Plus className="h-4 w-4 mr-2" />
                Add Task
              </Button>
            </div>
            <div className="text-sm text-muted-foreground">
              💡 Tip: Set priority and deadline to organize your tasks better. Default priority is medium.
            </div>
          </div>
        </div>

        {/* View Toggle and Controls */}
        <div className="flex flex-col lg:flex-row gap-4 mb-6">
          {/* View Toggle */}
          <div className="flex items-center gap-2">
            <Button
              variant={view === 'list' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setView('list')}
              className={view === 'list' 
                ? "bg-gradient-primary hover:opacity-90" 
                : "border-app-border hover:bg-app-surface-hover"
              }
            >
              <List className="h-4 w-4 mr-2" />
              List
            </Button>
            <Button
              variant={view === 'calendar' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setView('calendar')}
              className={view === 'calendar' 
                ? "bg-gradient-primary hover:opacity-90" 
                : "border-app-border hover:bg-app-surface-hover"
              }
            >
              <CalendarIcon className="h-4 w-4 mr-2" />
              Calendar
            </Button>
          </div>

          {/* Filters and Sorting */}
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="flex items-center gap-2">
              <Filter className="h-4 w-4 text-muted-foreground" />
              <Select value={filter} onValueChange={(value: FilterType) => setFilter(value)}>
                <SelectTrigger className="w-32 bg-card border-app-border hover:bg-app-surface-hover hover:border-primary transition-colors">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">All</SelectItem>
                  <SelectItem value="active" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">Active</SelectItem>
                  <SelectItem value="completed" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">Completed</SelectItem>
                  <SelectItem value="archived" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">📦 Archive</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="flex items-center gap-2">
              <span className="text-sm text-muted-foreground">Priority:</span>
              <Select value={priorityFilter} onValueChange={(value: PriorityType) => setPriorityFilter(value)}>
                <SelectTrigger className="w-32 bg-card border-app-border hover:bg-app-surface-hover hover:border-primary transition-colors">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">All</SelectItem>
                  <SelectItem value="high" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">🔴 High</SelectItem>
                  <SelectItem value="medium" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">🟡 Medium</SelectItem>
                  <SelectItem value="low" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">🟢 Low</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {view === 'list' && (
              <div className="flex items-center gap-2">
                <SortAsc className="h-4 w-4 text-muted-foreground" />
                <Select value={sort} onValueChange={(value: SortType) => setSort(value)}>
                  <SelectTrigger className="w-40 bg-card border-app-border hover:bg-app-surface-hover hover:border-primary transition-colors">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="newest" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">Newest First</SelectItem>
                    <SelectItem value="oldest" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">Oldest First</SelectItem>
                    <SelectItem value="alphabetical" className="hover:bg-app-surface-hover focus:bg-app-surface-hover cursor-pointer">Alphabetical</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            )}
          </div>
        </div>

        {/* Content based on view */}
        {view === 'list' ? (
          /* Tasks List */
          <div className="space-y-3">
            {sortedTasks.length === 0 ? (
              <div className="text-center py-12">
                <div className="text-6xl mb-4">📝</div>
                <h3 className="text-xl font-medium text-muted-foreground mb-2">
                  {filter === 'all' 
                    ? "No tasks yet" 
                    : `No ${filter} tasks`}
                </h3>
                <p className="text-muted-foreground">
                  {filter === 'all' 
                    ? "Add your first task to get started!"
                    : `You have no ${filter} tasks at the moment.`}
                </p>
              </div>
            ) : (
              sortedTasks.map((task, index) => (
                <div
                  key={task.id}
                  className="animate-fade-in"
                  style={{ animationDelay: `${index * 50}ms` }}
                >
                  <TaskItem
                    task={task}
                    onToggle={toggleTask}
                    onDelete={initiateDeleteTask}
                    onArchive={archiveTask}
                  />
                </div>
              ))
            )}
          </div>
        ) : (
          /* Calendar View */
          <WeekView
            tasks={sortedWeekTasks}
            onToggleTask={toggleTask}
            onDeleteTask={initiateDeleteTask}
            onArchiveTask={archiveTask}
            currentWeek={currentWeek}
            onWeekChange={setCurrentWeek}
          />
        )}
      </div>

      {/* Delete Confirmation Modal */}
      <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Подтвердите удаление</AlertDialogTitle>
            <AlertDialogDescription>
              Вы уверены, что хотите удалить задачу "{taskToDelete?.title}"? 
              Это действие нельзя будет отменить.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={cancelDeleteTask}>
              Отмена
            </AlertDialogCancel>
            <AlertDialogAction 
              onClick={confirmDeleteTask}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Да, удалить
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
};